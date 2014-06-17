package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/jackvalmadre/go-cv/detect"
	"github.com/jackvalmadre/go-cv/hog"
	"github.com/jackvalmadre/go-cv/imgpyr"
	"github.com/jackvalmadre/go-grideng/grideng"
	"github.com/jackvalmadre/go-ml"
	"github.com/nfnt/resize"
)

func init() {
	imgpyr.DefaultInterp = resize.Bilinear
	grideng.DefaultStdout = os.Stderr
}

func main() {
	var (
		sbin     = flag.Int("hog-sbin", 4, "Spatial bin parameter to HOG")
		pyrStep  = flag.Float64("pyr-step", 1.2, "Geometric scale steps in image pyramid")
		maxIOU   = flag.Float64("max-iou", 0, "Maximum intersection over union that two detections can have")
		margin   = flag.Int("margin", 0, "Spatial bin parameter to HOG")
		localMax = flag.Bool("local-max", true, "Suppress detections which are less than a neighbor?")
		minInter = flag.Float64("min-inter", 0.5, "Minimum intersection-over-union to validate a true positive")
	)

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, path.Base(os.Args[0]), "[flags] tmpl.gob inria/ roc.txt")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Tests a detector.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "inria/Test/")
		fmt.Fprintln(os.Stderr, "\tannotations/")
		fmt.Fprintln(os.Stderr, "\tannotations.lst")
		fmt.Fprintln(os.Stderr, "\tpos/")
		fmt.Fprintln(os.Stderr, "\tpos.lst")
		fmt.Fprintln(os.Stderr, "\tneg/")
		fmt.Fprintln(os.Stderr, "\tneg.lst")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
	}
	flag.Parse()
	grideng.ExecIfSlave()

	if flag.NArg() != 3 {
		flag.Usage()
		os.Exit(1)
	}
	var (
		tmplFile = flag.Arg(0)
		inriaDir = flag.Arg(1)
		rocFile  = flag.Arg(2)
	)

	detOpts := DetectOpts{
		HOGBin:   *sbin,
		PyrStep:  *pyrStep,
		MaxIOU:   *maxIOU,
		Margin:   *margin,
		LocalMax: *localMax,
	}

	// Load list of positive test image annotations.
	posAnnots, err := loadAnnots(path.Join(inriaDir, "Test", "annotations.lst"), inriaDir)
	if err != nil {
		log.Fatal(err)
	}
	// Load list of negative test images.
	negIms, err := loadLines(path.Join(inriaDir, "Test", "neg.lst"))
	if err != nil {
		log.Fatal(err)
	}

	// Load template.
	var tmpl *detect.FeatTmpl
	if err := loadGob(tmplFile, &tmpl); err != nil {
		log.Fatal(err)
	}
	log.Println("template size (pixels):", tmpl.Size)
	log.Println("template interior (pixels):", tmpl.Interior)
	log.Println("template size (features):", tmpl.Image.Size())
	if want, got := tmpl.Image.Size(), hog.FeatSize(tmpl.Size, hog.FGMRConfig(*sbin)); !got.Eq(want) {
		log.Fatalln("feature transform of patch is different size to weights")
	}

	// Test detector.
	annots := append(posAnnots, imsToAnnots(negIms)...)
	results, err := test(tmpl, annots, inriaDir, detOpts, *minInter)
	if err != nil {
		log.Fatal(err)
	}

	// Save results.
	enum := results.Enum()
	if err := saveResults(rocFile, enum); err != nil {
		log.Fatal(err)
	}
	fmt.Println("avgprec:", ml.ResultSet(enum).AveragePrecision())
}

func saveResults(fname string, results []ml.Result) error {
	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := writeResults(file, results); err != nil {
		return err
	}
	return nil
}

func writeResults(w io.Writer, results []ml.Result) error {
	avgprec := ml.ResultSet(results).AveragePrecision()
	if _, err := fmt.Fprintf(w, "# avg prec: %g\n", avgprec); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "TP\tTN\tFP\tFN"); err != nil {
		return err
	}
	for _, r := range results {
		s := fmt.Sprintf("%d\t%d\t%d\t%d", r.TP, r.TN, r.FP, r.FN)
		if _, err := fmt.Fprintln(w, s); err != nil {
			return err
		}
	}
	return nil
}

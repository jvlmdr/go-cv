package main

import (
	"github.com/jackvalmadre/go-cv/detect"
	"github.com/jackvalmadre/go-cv/featpyr"
	"github.com/jackvalmadre/go-cv/hog"
	"github.com/jackvalmadre/go-cv/imgpyr"
	"github.com/jackvalmadre/go-cv/rimg64"
	"github.com/jackvalmadre/go-ml"
	"github.com/nfnt/resize"

	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"path"
)

func main() {
	hogBin := flag.Int("hog-sbin", 4, "Spatial binning parameter to HOG")
	pyrStep := flag.Float64("pyr-step", 1.2, "Geometric steps in pyramid")
	valMinInter := flag.Float64("val-min-inter", 0.5, "Minimum intersection-over-union to validate a true positive")
	detMaxInter := flag.Float64("det-max-inter", 0.5, "Maximum permitted relative intersection before suppression")
	detLocalMax := flag.Bool("det-local-max", true, "Suppress detections which are less than a neighbor?")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, path.Base(os.Args[0]), "[flags] tmpl.gob inria/")
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
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}
	var (
		tmplFile = flag.Arg(0)
		inriaDir = flag.Arg(1)
	)
	opts := featpyr.DetectOpts{
		MaxInter: *detMaxInter,
		LocalMax: *detLocalMax,
	}
	imgpyr.DefaultInterp = resize.Bilinear

	// Load template.
	var tmpl *detect.FeatTmpl
	if err := loadGob(tmplFile, &tmpl); err != nil {
		log.Fatal(err)
	}

	// Evaluate template on each positive image and validate detections.
	pos, err := testPos(tmpl, inriaDir, *pyrStep, *hogBin, opts, *valMinInter)
	if err != nil {
		log.Fatal(err)
	}
	negDets, err := evalNeg(tmpl, inriaDir, *pyrStep, *hogBin, opts)
	if err != nil {
		log.Fatal(err)
	}
	neg := toFalsePos(negDets)
	results := pos.Merge(neg)

	if err := writeResults(os.Stdout, results.Enum()); err != nil {
		log.Fatal(err)
	}
}

func testPos(tmpl *detect.FeatTmpl, dir string, pyrStep float64, hogBin int, opts featpyr.DetectOpts, valMinInter float64) (*detect.ResultSet, error) {
	// Load list of annotations.
	anns, err := loadLines(path.Join(dir, "Test", "annotations.lst"))
	if err != nil {
		return nil, err
	}

	// Test each image and combine.
	var results *detect.ResultSet
	for i := range anns {
		imfile, refs, err := loadAnnotation(path.Join(dir, anns[i]))
		if err != nil {
			return nil, err
		}
		im, err := loadImage(path.Join(dir, imfile))
		if err != nil {
			return nil, err
		}
		// Get detections.
		dets := evalImage(tmpl, im, pyrStep, hogBin, opts)
		x := detect.ValidateMatch(dets, refs, valMinInter)
		results = results.Merge(x)
	}
	return results, nil
}

func evalNeg(tmpl *detect.FeatTmpl, dir string, pyrStep float64, hogBin int, opts featpyr.DetectOpts) ([]detect.Det, error) {
	// Load list of images.
	ims, err := loadLines(path.Join(dir, "Test", "neg.lst"))
	if err != nil {
		return nil, err
	}

	// Obtain detections from each image.
	var dets []detect.Det
	for i := range ims {
		im, err := loadImage(path.Join(dir, ims[i]))
		if err != nil {
			return nil, err
		}
		x := evalImage(tmpl, im, pyrStep, hogBin, opts)
		dets = detect.MergeDets(dets, x)
	}
	return dets, nil
}

func toFalsePos(dets []detect.Det) *detect.ResultSet {
	vals := make([]detect.ValDet, len(dets))
	for i, det := range dets {
		vals[i] = detect.ValDet{det, false}
	}
	return &detect.ResultSet{vals, 0}
}

func evalImage(tmpl *detect.FeatTmpl, im image.Image, pyrStep float64, hogBin int, opts featpyr.DetectOpts) []detect.Det {
	// Construct image pyramid.
	scales := imgpyr.Scales(im.Bounds().Size(), tmpl.PixSize(), pyrStep)
	pixpyr := imgpyr.New(im, scales)
	// Construct HOG pyramid.
	fn := func(rgb *rimg64.Multi) *rimg64.Multi { return hog.FGMR(rgb, hogBin) }
	pyr := featpyr.New(pixpyr, fn, hogBin)
	// Search feature pyramid.
	dets := featpyr.Detect(pyr, tmpl, opts)
	return dets
}

func loadImage(fname string) (image.Image, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	im, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return im, nil
}

func loadAnnotation(fname string) (im string, refs []image.Rectangle, err error) {
	lines, err := loadLines(fname)
	if err != nil {
		return "", nil, err
	}

	// Get image file.
	im = func() string {
		var s string
		for _, line := range lines {
			n, _ := fmt.Sscanf(line, "Image filename : %q", &s)
			if n == 1 {
				return s
			}
		}
		return ""
	}()

	// Get object bounds.
	refs = func() []image.Rectangle {
		var rs []image.Rectangle
		for _, line := range lines {
			var class string
			var xmin, ymin, xmax, ymax int
			const format = "Bounding box for object %d %q (Xmin, Ymin) - (Xmax, Ymax) : (%d, %d) - (%d, %d)"
			n, _ := fmt.Sscanf(line, format, new(int), &class, &xmin, &ymin, &xmax, &ymax)
			if n != 6 {
				continue
			}
			if class != "PASperson" {
				panic("found non-person: " + class)
			}
			rs = append(rs, image.Rect(xmin, ymin, xmax, ymax))
		}
		return rs
	}()
	return im, refs, nil
}

func loadLines(fname string) ([]string, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var x []string
	s := bufio.NewScanner(file)
	for s.Scan() {
		x = append(x, s.Text())
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return x, nil
}

func loadGob(fname string, x interface{}) error {
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()
	return gob.NewDecoder(file).Decode(x)
}

func writeResults(w io.Writer, results []ml.Result) error {
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

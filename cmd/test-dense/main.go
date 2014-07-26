package main

import (
	"github.com/jvlmdr/go-cv/detect"
	"github.com/jvlmdr/go-cv/hog"
	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-cv/slide"
	"github.com/jvlmdr/go-ml"

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
	"sort"
)

func main() {
	sbin := flag.Int("sbin", 4, "Spatial binning parameter to HOG")
	posDir := flag.String("pos-dir", "", "Directory of the positive images")
	negDir := flag.String("neg-dir", "", "Directory of the negative images")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, path.Base(os.Args[0]), "[flags] tmpl.gob pos.txt neg.txt")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Tests a detector.")
		fmt.Fprintln(os.Stderr, "Positive images must be cropped to same size as template.")
		fmt.Fprintln(os.Stderr, "Negative images are evaluated densely without non-maximum suppression.")
		fmt.Fprintln(os.Stderr, "No pyramids are used.")
		fmt.Fprintln(os.Stderr, "Results for varying threshold are printed to stdout.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
	}
	flag.Parse()
	if flag.NArg() != 3 {
		flag.Usage()
		os.Exit(1)
	}

	var (
		tmplFile = flag.Arg(0)
		posFile  = flag.Arg(1)
		negFile  = flag.Arg(2)
	)

	// Load template.
	var tmpl *detect.FeatTmpl
	if err := loadGob(tmplFile, &tmpl); err != nil {
		log.Fatal(err)
	}

	// Evaluate detector on all images in the positive set.
	posScores, err := evalExamplesFile(tmpl.Image, posFile, *posDir, *sbin)
	if err != nil {
		log.Fatal(err)
	}
	// Evaluate detector on all images in the negative set.
	negScores, err := evalImagesFile(tmpl.Image, negFile, *negDir, *sbin)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("sort positives: ", len(posScores))
	log.Print("sort negatives: ", len(negScores))

	// TP, FP, TN, FN data for ROC or analogous curve.
	results := ml.EnumerateResults(posScores, negScores)
	//	results = summarize(results)

	if err := writeResults(os.Stdout, results); err != nil {
		log.Fatal(err)
	}
}

func evalExamplesFile(tmpl *rimg64.Multi, imsFile string, dir string, sbin int) ([]float64, error) {
	ims, err := loadLines(imsFile)
	if err != nil {
		return nil, err
	}
	return evalExamples(tmpl, ims, dir, sbin)
}

func evalExamples(tmpl *rimg64.Multi, ims []string, dir string, sbin int) ([]float64, error) {
	var scores []float64
	for i, file := range ims {
		log.Printf("pos: %d/%d: %s", i+1, len(ims), file)
		im, err := loadFeatImage(path.Join(dir, file), sbin)
		if err != nil {
			return nil, err
		}
		// Check that dimensions match.
		if im.Width != tmpl.Width || im.Height != tmpl.Height {
			err := fmt.Errorf(
				"different size: template %dx%d, image %dx%d",
				tmpl.Width, tmpl.Height, im.Width, im.Height,
			)
			return nil, err
		}
		scores = append(scores, dot(tmpl, im))
	}
	return scores, nil
}

func evalImagesFile(tmpl *rimg64.Multi, imsFile string, dir string, sbin int) ([]float64, error) {
	ims, err := loadLines(imsFile)
	if err != nil {
		return nil, err
	}
	return evalImages(tmpl, ims, dir, sbin)
}

func evalImages(tmpl *rimg64.Multi, ims []string, dir string, sbin int) ([]float64, error) {
	var scores []float64
	for i, file := range ims {
		log.Printf("neg: %d/%d: %s", i+1, len(ims), file)
		im, err := loadFeatImage(path.Join(dir, file), sbin)
		if err != nil {
			return nil, err
		}

		// Obtain response to sliding window.
		resp := slide.CorrMulti(im, tmpl)
		scores = append(scores, resp.Elems...)
	}
	return scores, nil
}

func dot(x, y *rimg64.Multi) float64 {
	var d float64
	for i := 0; i < x.Width; i++ {
		for j := 0; j < x.Height; j++ {
			for k := 0; k < x.Channels; k++ {
				d += x.At(i, j, k) * y.At(i, j, k)
			}
		}
	}
	return d
}

func loadFeatImage(fname string, sbin int) (*rimg64.Multi, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	im, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return hog.HOG(rimg64.FromColor(im), hog.FGMRConfig(sbin)), nil
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

func saveGob(fname string, x interface{}) error {
	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer file.Close()
	return gob.NewEncoder(file).Encode(x)
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

func summarize(in []ml.Result) []ml.Result {
	var out []ml.Result
	for i, r := range in {
		if i == 0 {
			// Append first.
			out = append(out, r)
		} else if i < len(in)-1 {
			if (in[i-1].TP == in[i].TP) != (in[i].TP == in[i+1].TP) {
				out = append(out, r)
			}
		} else {
			// Append last.
			out = append(out, r)
		}
	}
	return out
}

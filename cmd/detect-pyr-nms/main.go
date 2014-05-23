package main

import (
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
	"path"

	"github.com/jackvalmadre/go-cv/detect"
	"github.com/jackvalmadre/go-cv/feat"
	"github.com/jackvalmadre/go-cv/featpyr"
	"github.com/jackvalmadre/go-cv/hog"
	"github.com/jackvalmadre/go-cv/imgpyr"
	"github.com/jackvalmadre/go-cv/imsamp"
	"github.com/nfnt/resize"
)

func init() {
	imgpyr.DefaultInterp = resize.Bilinear
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, path.Base(os.Args[0]), "[flags] tmpl.gob image.(jpg|png) detections.json")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Runs a detector on an image with non-maximum suppression.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Options:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}

func main() {
	var (
		sbin     = flag.Int("sbin", 4, "Spatial binning parameter to HOG")
		margin   = flag.Int("margin", 0, "Margin to add around images before computing features")
		step     = flag.Float64("pyr-step", 1.1, "Geometric step to use in image pyramid")
		maxinter = flag.Float64("max-intersect", 0.5, "Maximum overlap of detections. Zero means detections can't overlap at all, one means they can overlap entirely.")
		localmax = flag.Bool("local-max", true, "Detections cannot score less than a neighbor")
	)

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 3 {
		flag.Usage()
		os.Exit(1)
	}
	var (
		tmplFile = flag.Arg(0)
		imFile   = flag.Arg(1)
		detsFile = flag.Arg(2)
	)

	// Load image.
	im, err := loadImage(imFile)
	if err != nil {
		log.Fatal(err)
	}
	// Construct pyramid.
	scales := imgpyr.Scales(im.Bounds().Size(), image.Pt(24, 24), *step)
	phi := hog.Transform{hog.FGMRConfig(*sbin)}
	pad := feat.Pad{feat.Margin{*margin, *margin, *margin, *margin}, imsamp.Continue}
	pyr := featpyr.NewPad(imgpyr.New(im, scales), phi, pad)

	// Load template.
	var tmpl *detect.FeatTmpl
	if err := loadGob(tmplFile, &tmpl); err != nil {
		log.Fatal(err)
	}

	detopts := detect.DetFilter{
		LocalMax: *localmax,
		MinScore: math.Inf(-1),
	}
	// Use intersection-over-union criteria for non-max suppression.
	overlap := func(a, b image.Rectangle) bool {
		return detect.IOU(a, b) > *maxinter
	}
	suppropts := detect.SupprFilter{
		MaxNum:  0,
		Overlap: overlap,
	}
	dets := detect.Pyramid(pyr, tmpl, detopts, suppropts)

	if err := saveJSON(detsFile, dets); err != nil {
		log.Fatal(err)
	}

	for i, det := range dets {
		r := det.Rect
		cmd := fmt.Sprintf("rectangle %d,%d %d,%d", r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)
		fmt.Printf("convert %s -fill none -stroke white -draw '%s' det_%06d.jpg\n", imFile, cmd, i)
	}
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

func loadGob(fname string, x interface{}) error {
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()
	return gob.NewDecoder(file).Decode(x)
}

func saveJSON(fname string, x interface{}) error {
	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(x)
}

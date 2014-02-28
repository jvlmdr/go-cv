package main

import (
	"github.com/jackvalmadre/go-cv/featpyr"
	"github.com/jackvalmadre/go-cv/hog"
	"github.com/jackvalmadre/go-cv/imgpyr"
	"github.com/jackvalmadre/go-cv/rimg64"
	"github.com/nfnt/resize"

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
)

func main() {
	sbin := flag.Int("sbin", 4, "Spatial binning parameter to HOG")
	step := flag.Float64("pyr-step", 1.2, "Geometric step to use in image pyramid")
	maxinter := flag.Float64("max-intersect", 0.5, "Maximum overlap of detections. Zero means detections can't overlap at all, one means they can overlap entirely.")
	localmax := flag.Bool("local-max", true, "Detections cannot score less than a neighbor")
	width := flag.Int("tmpl-pix-width", 0, "Width of template in pixels")
	height := flag.Int("tmpl-pix-height", 0, "Height of template in pixels")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, path.Base(os.Args[0]), "[flags] tmpl.gob image.(jpg|png) detections.json")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Runs a detector on an image with non-maximum suppression.")
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
		imFile   = flag.Arg(1)
		detsFile = flag.Arg(2)
	)

	if *width <= 0 || *height <= 0 {
		log.Fatal("width and height must be positive")
	}
	pixsize := image.Pt(*width, *height)

	// Load image.
	im, err := loadImage(imFile)
	if err != nil {
		log.Fatal(err)
	}
	// Construct pyramid.
	fn := func(x *rimg64.Multi) *rimg64.Multi {
		return hog.HOG(x, FGMRConfig(*sbin))
	}
	scales := imgpyr.Scales(im.Bounds().Size(), image.Pt(64, 64), *step)
	pyr := featpyr.New(im, scales, fn, *sbin)

	// Load template.
	var tmpl *rimg64.Multi
	if err := loadGob(tmplFile, &tmpl); err != nil {
		log.Fatal(err)
	}

	imgpyr.DefaultInterp = resize.Bilinear
	opts := featpyr.DetectOpts{MaxInter: *maxinter, MinScore: math.Inf(-1), LocalMax: *localmax}
	dets := featpyr.Detect(pyr, tmpl, pixsize, opts)

	if err := saveJSON(detsFile, dets); err != nil {
		log.Fatal(err)
	}

	//	for i, det := range dets {
	//		r := det.Pos
	//		cmd := fmt.Sprintf("rectangle %d,%d %d,%d", r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)
	//		fmt.Printf("convert %s -fill none -stroke white -draw '%s' det_%06d.jpg\n", imFile, cmd, i)
	//	}
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

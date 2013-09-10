package main

import (
	"github.com/jackvalmadre/go-cv"
	"github.com/jackvalmadre/go-cv/hog"
	"github.com/jackvalmadre/go-cv/imgpyr"
	"github.com/jackvalmadre/go-cv/slide"
	"github.com/nfnt/resize"

	"encoding/gob"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
)

type PyrPos struct {
	Level int
	Pos   image.Point
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage of %s:\n", os.Args[0])
	fmt.Fprintln(os.Stderr, os.Args[0], "image.(jpg|png) weights.gob report.html")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Finds multi-scale detections without thresholding")
	fmt.Fprintln(os.Stderr)
}

var (
	binSize int
	geoStep float64
)

func init() {
	flag.IntVar(&binSize, "bin-size", 4, "HOG bin size")
	flag.Float64Var(&geoStep, "step", 1.1, "Geometric step in pyramid")
}

func main() {
	flag.Usage = usage
	flag.Parse()

	// Read other arguments.
	if flag.NArg() != 3 {
		flag.Usage()
		os.Exit(1)
	}
	var (
		imgFile     = flag.Arg(0)
		weightsFile = flag.Arg(1)
		reportFile  = flag.Arg(2)
	)

	if geoStep <= 1 {
		log.Fatalln("geometric step must be greater than 1")
	}

	// Load image from file.
	img, err := loadImage(imgFile)
	if err != nil {
		log.Fatalln("could not load image:", err)
	}
	log.Println("image bounds:", img.Bounds())
	// Load weights from file.
	weights, err := loadWeights(weightsFile)
	if err != nil {
		log.Fatalln("could not load weights:", err)
	}

	log.Println("transform original image")
	// Compute feature transform of original image.
	featImg := hogImage(img, binSize)
	// Check image and detector sizes.
	if weights.Width > featImg.Width || weights.Height > featImg.Height {
		log.Fatalln("detector larger than feature image")
	}
	// Plan pyramid levels.
	minXScale := float64(weights.Width) / float64(featImg.Width)
	minYScale := float64(weights.Height) / float64(featImg.Height)
	minScale := math.Max(minXScale, minYScale)
	// Construct pyramid.
	seq := imgpyr.Sequence(1, 1/geoStep, minScale)
	log.Println("number of zoom levels:", seq.Len)
	log.Println("compute pyramid")
	pyr := imgpyr.NewInterp(img, seq, resize.NearestNeighbor)

	// Transform every level of the pyramid.
	log.Println("transform each level of pyramid")
	featImgs := make([]cv.RealVectorImage, len(pyr.Levels))
	for i, lvl := range pyr.Levels {
		featImgs[i] = hogImage(lvl, binSize)
	}

	// Run detector on every level of the pyramid.
	log.Println("search each level of pyramid")
	scores := make([]cv.RealImage, len(pyr.Levels))
	for i, featImg := range featImgs {
		scores[i] = slide.CorrelateVectorImages(featImg, weights)
	}

	// Get detections.
	log.Println("non-maxima suppression")
	dets := nonMaxSupp(scores, pyr.Scales, weights.ImageSize(), 0.5)

	// Generate report.
	report(reportFile, pyr, dets, weights.ImageSize(), scores, binSize)
}

type Det struct {
	Rect  image.Rectangle
	Score float64
}

// Separate function to ensure that file is not open longer than needed.
func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
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

// Separate function to ensure that file is not open longer than needed.
func loadWeights(filename string) (cv.RealVectorImage, error) {
	file, err := os.Open(filename)
	if err != nil {
		return cv.RealVectorImage{}, err
	}
	defer file.Close()

	var detector struct {
		Weights cv.RealVectorImage
	}
	if err := gob.NewDecoder(file).Decode(&detector); err != nil {
		return cv.RealVectorImage{}, err
	}
	return detector.Weights, nil
}

// Returns HOG image of a visual image.
func hogImage(img image.Image, binSize int) cv.RealVectorImage {
	return hog.HOG(cv.ColorImageToReal(img), binSize)
}

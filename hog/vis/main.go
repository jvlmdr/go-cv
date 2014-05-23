package main

import (
	"github.com/jackvalmadre/go-cv"
	"github.com/jackvalmadre/go-cv/hog"

	"flag"
	"fmt"
	"log"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s input.gob output.png\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	log.SetOutput(os.Stdout)

	var (
		name string
		cell int
	)
	flag.StringVar(&name, "weights", "signed", `in {"signed", "pos", "neg", "abs"}`)
	flag.IntVar(&cell, "cell", 32, "Size to render cells (pixels)")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}
	input := flag.Arg(0)
	output := flag.Arg(1)

	fmt.Println("Opening image...")
	var feat cv.RealVectorImage
	if err := decode(input, &feat); err != nil {
		log.Fatal(err)
	}

	var weights hog.WeightSet
	switch name {
	case "pos":
		weights = hog.Pos
	case "neg":
		weights = hog.Neg
	case "abs":
		weights = hog.Abs
	default:
		weights = hog.Signed
	}

	fmt.Println("Rendering visualization...")
	pic := hog.Vis(feat, weights, cell)

	fmt.Println("Saving image...")
	if err := saveImage(output, pic); err != nil {
		fmt.Println(err)
		return
	}
}

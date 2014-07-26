package main

import (
	"github.com/jvlmdr/go-cv/imgpyr"

	"flag"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage of %s:\n", os.Args[0])
	fmt.Fprintln(os.Stderr, os.Args[0], "image pyramid")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	var (
		imgFile = flag.Arg(0)
		visFile = flag.Arg(1)
	)

	img, err := loadImage(imgFile)
	if err != nil {
		log.Fatalln(err)
	}

	const r = 1.5
	seq := imgpyr.Sequence(1, 1/r, 0.2)
	pyr := imgpyr.New(img, seq)

	vis := pyr.Visualize()
	if err := saveImage(visFile, vis); err != nil {
		log.Fatalln(err)
	}
}

func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func saveImage(filename string, img image.Image) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := jpeg.Encode(file, img, nil); err != nil {
		return err
	}
	return nil
}

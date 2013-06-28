package main

import (
	"flag"
	"fmt"
	"github.com/jackvalmadre/go-cv"
	"log"
	"os"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage of", os.Args[0])
	fmt.Fprintln(os.Stderr, os.Args[0], "im.gob")
	flag.PrintDefaults()
}

func main() {
	log.SetOutput(os.Stdout)
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		return
	}
	filename := flag.Arg(0)

	var im cv.RealVectorImage
	if err := decodeGob(filename, &im); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%dx%d image with %d channels\n", im.Width, im.Height, im.Channels)
}

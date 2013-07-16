package main

import (
	"flag"
	"fmt"
	"github.com/jackvalmadre/go-cv"
	"github.com/jackvalmadre/lin-go/vec"
	"image"
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
	flag.Usage = usage
	signed := flag.Bool("signed", false, "Treat pixels as signed")
	negate := flag.Bool("negate", false, "Negate image")
	cellSize := flag.Int("cell-size", 32, "Size to render cells (pixels)")
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}
	input := flag.Arg(0)
	output := flag.Arg(1)

	fmt.Println("Opening image...")
	var phi cv.RealVectorImage
	if err := decode(input, &phi); err != nil {
		log.Fatal(err)
	}

	if *negate {
		vec.CopyTo(phi.Vec(), vec.Scale(-1, phi.Vec()))
	}

	fmt.Println("Rendering visualization...")
	var pic image.Image
	if *signed {
		pic = cv.SignedHOGImage(phi, *cellSize)
	} else {
		pic = cv.HOGImage(phi, *cellSize)
	}

	fmt.Println("Saving image...")
	if err := saveImage(output, pic); err != nil {
		fmt.Println(err)
		return
	}
}

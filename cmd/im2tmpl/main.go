package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"image"
	"os"
	"path"

	"github.com/jvlmdr/go-cv/detect"
	"github.com/jvlmdr/go-cv/rimg64"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s weights.(gob|csv)\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var (
		width  = flag.Int("width", 0, "Width in pixels")
		height = flag.Int("height", 0, "Height in pixels")
		top    = flag.Int("top", 0, "Inset from pixel size to give interior")
		left   = flag.Int("left", 0, "Inset from pixel size to give interior")
		bottom = flag.Int("bottom", 0, "Inset from pixel size to give interior")
		right  = flag.Int("right", 0, "Inset from pixel size to give interior")
	)

	flag.Parse()
	flag.Usage = usage
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}
	imname := flag.Arg(0)
	tmplname := flag.Arg(1)

	imfile, err := os.Open(imname)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer imfile.Close()

	var weights *rimg64.Multi
	if path.Ext(imname) == ".csv" {
		var err error
		weights, err = readImageCSV(imfile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		if err := gob.NewDecoder(imfile).Decode(&weights); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	tmpl := &detect.FeatTmpl{
		weights,
		image.Pt(*width, *height),
		image.Rect(*left, *top, *width-*right, *height-*bottom),
	}

	tmplfile, err := os.Create(tmplname)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer tmplfile.Close()

	if err := gob.NewEncoder(tmplfile).Encode(tmpl); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

package main

import (
	"github.com/jackvalmadre/go-cv/hog"
	"github.com/jackvalmadre/go-cv/rimg64"

	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func loadFeatureImage(file string, sbin int) (*rimg64.Multi, error) {
	f, err := loadRealImage(file)
	if err != nil {
		return nil, err
	}
	g := hog.HOG(f, hog.FGMRConfig(sbin))
	return g, nil
}

func loadRealImage(file string) (*rimg64.Multi, error) {
	im, err := loadImage(file)
	if err != nil {
		return nil, err
	}
	f := rimg64.FromColor(im)
	return f, nil
}

func loadImage(name string) (image.Image, error) {
	file, err := os.Open(name)
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

package main

import (
	"encoding/gob"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
)

func decode(name string, e interface{}) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := gob.NewDecoder(file)
	err = dec.Decode(e)
	if err != nil {
		return err
	}

	return nil
}

func encode(name string, e interface{}) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	err = enc.Encode(e)
	if err != nil {
		return err
	}

	return nil
}

func loadImage(name string) (*image.Image, error) {
	// Load image.
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// Decode to bitmap.
	im, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return &im, nil
}

func saveImage(name string, im image.Image) error {
	// Load image.
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	// Decode to bitmap.
	err = png.Encode(file, im)
	if err != nil {
		return err
	}
	return nil
}

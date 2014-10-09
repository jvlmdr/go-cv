package feat

import (
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
)

func init() {
	RegisterImage("gray", func() ImageSpec { return &SimpleImageSpec{NewGray()} })
	RegisterImage("rgb", func() ImageSpec { return &SimpleImageSpec{NewRGB()} })
}

// NewGray returns a transform which converts the image to a real-valued gray image.
func NewGray() ImageMarshalable {
	return newImageFunc(toGray, 1, "gray")
}

// toGray never encounters an error.
func toGray(im image.Image) (*rimg64.Multi, error) {
	x := rimg64.FromGray(im)
	// Convert into multi-channel image with one channel.
	y := &rimg64.Multi{x.Elems, x.Width, x.Height, 1}
	return y, nil
}

func toReal(im image.Image) (*rimg64.Multi, error) {
	return rimg64.FromColor(im), nil
}

// NewRGB returns a transform which converts the image to a real-valued image.
// The channels are red, green and blue.
func NewRGB() ImageMarshalable {
	return &marshalableImageFunc{imageFunc: imageFunc{apply: toReal, rate: 1}, name: "rgb"}
}

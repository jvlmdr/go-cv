package featset

import (
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
)

func init() {
	RegisterImage("gray", func() Image { return new(Gray) })
	RegisterImage("rgb", func() Image { return new(RGB) })
}

type Gray struct{}

func (phi *Gray) Rate() int { return 1 }

func (phi *Gray) Apply(im image.Image) (*rimg64.Multi, error) {
	return toGray(im), nil
}

func (phi *Gray) Marshaler() *ImageMarshaler {
	return &ImageMarshaler{"gray", nil}
}

func (phi *Gray) Transform() Image { return phi }

type RGB struct{}

func (phi *RGB) Rate() int { return 1 }

func (phi *RGB) Apply(im image.Image) (*rimg64.Multi, error) {
	return rimg64.FromColor(im), nil
}

func (phi *RGB) Marshaler() *ImageMarshaler {
	return &ImageMarshaler{"rgb", nil}
}

func (phi *RGB) Transform() Image {
	return phi
}

// toGray never encounters an error.
func toGray(im image.Image) *rimg64.Multi {
	x := rimg64.FromGray(im)
	// Convert into multi-channel image with one channel.
	y := &rimg64.Multi{x.Elems, x.Width, x.Height, 1}
	return y
}

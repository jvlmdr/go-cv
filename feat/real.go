package feat

import (
	"image"

	//"github.com/gonum/floats"
	"github.com/jvlmdr/go-cv/rimg64"
)

func init() {
	RegisterImage("gray", NewGray)
	RegisterImage("color", NewRGB)
}

// NewGray returns a transform which converts the image to a real-valued gray image.
func NewGray() Image {
	return NewImageFunc(toGray, 1)
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
func NewRGB() Image {
	return NewImageFunc(toReal, 1)
}

//	// NewGrayNorm converts an image to gray and normalizes it
//	// to have zero mean and squared norm equal to its area.
//	// Note that this only satisfies the Image spec
//	// in the limit and after re-normalization.
//	func NewGrayNorm() Image {
//		f := func(im image.Image) (*rimg64.Multi, error) {
//			x := rimg64.FromGray(im)
//			normalize(x.Elems, float64(x.Width*x.Height))
//			y := &rimg64.Multi{x.Elems, x.Width, x.Height, 1}
//			return y, nil
//		}
//		return NewImageFunc(f, 1)
//	}
//
//	// Modifies slice to have zero mean and square norm as given.
//	func normalize(x []float64, sqrnorm float64) {
//		// Subtract mean (in-place).
//		mean := floats.Sum(x) / float64(len(x))
//		floats.AddConst(-mean, x)
//		// Set squared-norm equal (in-place).
//		floats.Scale(sqrnorm/floats.Dot(x, x), x)
//	}

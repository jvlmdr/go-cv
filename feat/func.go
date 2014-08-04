package feat

import (
	"image"

	"github.com/gonum/floats"
	"github.com/jvlmdr/go-cv/rimg64"
)

func init() {
	Register("gray", Gray)
	Register("color", Color)
}

type ApplyFunc func(image.Image) *rimg64.Multi

func Func(f ApplyFunc, rate int) Transform {
	return funcTransform{f, rate}
}

type funcTransform struct {
	apply ApplyFunc
	rate  int
}

func (t funcTransform) Apply(im image.Image) *rimg64.Multi {
	return t.apply(im)
}

func (t funcTransform) Rate() int {
	return t.rate
}

// Gray returns a transform which convert the image to gray.
func Gray() Transform {
	f := func(im image.Image) *rimg64.Multi {
		g := rimg64.FromGray(im)
		// Convert into multi-channel image with one channel.
		return &rimg64.Multi{g.Elems, g.Width, g.Height, 1}
	}
	return Func(f, 1)
}

func Color() Transform {
	return Func(rimg64.FromColor, 1)
}

// GrayNorm converts an image to gray and normalizes it
// to have zero mean and squared norm equal to its area.
// Note that this only satisfies the Transform spec
// in the limit and after re-normalization.
func GrayNorm() Transform {
	f := func(im image.Image) *rimg64.Multi {
		g := rimg64.FromGray(im)
		normalize(g.Elems, float64(g.Width*g.Height))
		return &rimg64.Multi{g.Elems, g.Width, g.Height, 1}
	}
	return Func(f, 1)
}

/*
// ColorNorm normalizes a color image
// to have zero mean and squared norm equal to three times its area.
// Channels are normalized jointly.
// Note that this only satisfies the Transform spec
// in the limit and after re-normalization.
func ColorNorm() Transform {
	f := func(im image.Image) *rimg64.Multi {
		g := rimg64.FromColor(im)
		normalize(g.Elems, float64(g.Width*g.Height*g.Channels))
		return g
	}
	return Func(f, 1)
}
*/

// Modifies slice to have zero mean and square norm as given.
func normalize(x []float64, sqrnorm float64) {
	// Subtract mean (in-place).
	mean := floats.Sum(x) / float64(len(x))
	floats.AddConst(-mean, x)
	// Set squared-norm equal (in-place).
	floats.Scale(sqrnorm/floats.Dot(x, x), x)
}

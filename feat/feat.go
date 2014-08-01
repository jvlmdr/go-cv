package feat

import (
	"image"

	"github.com/jvlmdr/go-cv/imsamp"
	"github.com/jvlmdr/go-cv/rimg64"
)

// Transform defines a common interface for feature transforms.
//
// Feature transforms are assumed to have a positive integer downsample rate.
// If the input image f(x, y) has size (m, n) with domain
// 	x = 0, ..., m - 1
// 	y = 0, ..., n - 1
// and produces a feature image g(u, v) of size (p, q) with domain
// 	u = 0, ..., p - 1
// 	v = 0, ..., q - 1
// then calling Apply() on any sub-image
// 	x = rate*left, ..., m - 1 - rate*right
// 	y = rate*top, ..., n - 1 - rate*bottom
// must produce the feature image
// 	u = left, ..., p - 1 - right
// 	v = top, ..., q - 1 - bottom
// where (left, right, top, bottom) are non-negative integers describing an inset on each side.
type Transform interface {
	// Function to compute transform on image.
	Apply(im image.Image) *rimg64.Multi
	// Integer downsample rate.
	Rate() int
}

type Margin struct {
	Top, Left, Bottom, Right int
}

func (m Margin) Empty() bool {
	return m == Margin{}
}

func (m Margin) TopLeft() image.Point {
	return image.Pt(m.Left, m.Top)
}

func (m Margin) AddTo(r image.Rectangle) image.Rectangle {
	r.Min.X -= m.Left
	r.Min.Y -= m.Top
	r.Max.X += m.Right
	r.Max.Y += m.Bottom
	return r
}

// ApplyPad pads the image before computing the feature transform.
func ApplyPad(t Transform, im image.Image, pad Pad) *rimg64.Multi {
	if pad.Margin.Empty() {
		// Take a shortcut.
		return t.Apply(im)
	}
	im = imsamp.Rect(im, pad.Margin.AddTo(im.Bounds()), pad.Extend)
	return t.Apply(im)
}

type Pad struct {
	Margin Margin
	Extend imsamp.At
}

func NoPad() Pad {
	return Pad{}
}

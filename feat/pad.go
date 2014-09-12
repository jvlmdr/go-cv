package feat

import (
	"image"

	"github.com/jvlmdr/go-cv/imsamp"
	"github.com/jvlmdr/go-cv/rimg64"
)

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

func UniformMargin(x int) Margin {
	return Margin{x, x, x, x}
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

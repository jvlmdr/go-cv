package feat

import (
	"image"

	"github.com/jackvalmadre/go-cv/imsamp"
	"github.com/jackvalmadre/go-cv/rimg64"
)

// Describes a feature transform.
type Transform interface {
	// Function to compute transform on image.
	Apply(im image.Image) *rimg64.Multi
	// Integer down-sample rate.
	Rate() int
}

type Margin struct {
	Top, Left, Bottom, Right int
}

func (m Margin) Empty() bool {
	return m.Top == 0 && m.Left == 0 && m.Bottom == 0 && m.Right == 0
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

// Extends the image by the given margin before computing the feature transform.
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

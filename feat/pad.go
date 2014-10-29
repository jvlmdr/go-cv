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

func (m Margin) BottomRight() image.Point {
	return image.Pt(m.Right, m.Bottom)
}

func (m Margin) AddTo(r image.Rectangle) image.Rectangle {
	return image.Rectangle{r.Min.Sub(m.TopLeft()), r.Max.Add(m.BottomRight())}
}

func (m Margin) SubFrom(r image.Rectangle) image.Rectangle {
	return image.Rectangle{r.Min.Add(m.TopLeft()), r.Max.Sub(m.BottomRight())}
}

func UniformMargin(x int) Margin {
	return Margin{x, x, x, x}
}

// ApplyPad pads the image before computing the feature transform.
func ApplyPad(t Image, im image.Image, pad Pad) (*rimg64.Multi, error) {
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

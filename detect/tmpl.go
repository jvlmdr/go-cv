package detect

import (
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
)

type FeatTmpl struct {
	// Template in feature space.
	Image *rimg64.Multi
	Bias  float64
	// Size in pixels.
	Size image.Point
	// Interior of window in pixels.
	Interior image.Rectangle
}

func (tmpl *FeatTmpl) Bounds() image.Rectangle {
	return image.Rectangle{image.ZP, tmpl.Size}
}

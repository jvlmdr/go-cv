package detect

import (
	"github.com/jackvalmadre/go-cv/rimg64"

	"image"
)

// Feature template.
type FeatTmpl struct {
	// Template in feature space.
	Image *rimg64.Multi
	// Size in pixels.
	Size image.Point
	// Interior of window in pixels.
	Interior image.Rectangle
}

func (tmpl *FeatTmpl) Bounds() image.Rectangle {
	return image.Rectangle{image.ZP, tmpl.Size}
}

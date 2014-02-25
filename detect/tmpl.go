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
	PixWidth  int
	PixHeight int
}

func (tmpl *FeatTmpl) PixSize() image.Point {
	return image.Pt(tmpl.PixWidth, tmpl.PixHeight)
}

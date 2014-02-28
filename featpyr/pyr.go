package featpyr

import (
	"github.com/jackvalmadre/go-cv/imgpyr"
	"github.com/jackvalmadre/go-cv/rimg64"

	"image"
	"log"
)

// Level zero is the original image.
// Level i has size (width, height) * Scales.At(i).
type Pyramid struct {
	Images *imgpyr.Pyramid
	// Feature transform of each level in pyramid.
	Feats []*rimg64.Multi
	// Integer downsample rate of features.
	Rate int
}

// Retrieves the pixel image.
// Level zero is the original image.
func (pyr *Pyramid) Image(level int) image.Image {
	return pyr.Images.Levels[level]
}

// Accesses the scale of level i.
func (pyr *Pyramid) Scale(i int) float64 {
	return pyr.Images.Scales.At(i)
}

func (pyr *Pyramid) ToImagePoint(pt imgpyr.Point) image.Point {
	return pointAt(pt, pyr.Images.Scales, pyr.Rate)
}

// Given point in feature pyramid and rectangle in pixel co-ords.
func (pyr *Pyramid) ToImageRect(min imgpyr.Point, interior image.Rectangle) image.Rectangle {
	return rectAt(min, pyr.Images.Scales, pyr.Rate, interior)
}

// Reserve name "Feat" in case this eventually becomes a struct.

// Transforms a floating-point color or gray image into a feature image.
type FeatFunc func(x *rimg64.Multi) *rimg64.Multi

// Constructs a feature pyramid.
// Level i has dimension (width, height) * scales[i].
func New(images *imgpyr.Pyramid, fn FeatFunc, rate int) *Pyramid {
	feats := make([]*rimg64.Multi, len(images.Levels))
	for i, x := range images.Levels {
		feats[i] = fn(rimg64.FromColor(x))
	}
	log.Print("finished computing feature transform: ", len(images.Levels))
	return &Pyramid{images, feats, rate}
}

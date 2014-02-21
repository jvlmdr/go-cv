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
	Feats  []*rimg64.Multi
	Rate   int
}

// Retrieves the pixel image.
// Level zero is the original image.
func (pyr *Pyramid) Image(level int) image.Image {
	return pyr.Images.Levels[level]
}

// Retrieves the feature image.
// Level zero is the original image.
func (pyr *Pyramid) Feat(level int) *rimg64.Multi {
	return pyr.Feats[level]
}

// Accesses the scale of level i.
func (pyr *Pyramid) Scale(i int) float64 {
	return pyr.Images.Scales.At(i)
}

// Reserve name "Feat" in case this eventually becomes a struct.

// Transforms a floating-point color or gray image into a feature image.
type FeatFunc func(x *rimg64.Multi) *rimg64.Multi

// Constructs a feature pyramid.
// Level i has dimension (width, height) * scales[i].
func New(im image.Image, scales imgpyr.GeoSeq, fn FeatFunc, rate int) *Pyramid {
	images := imgpyr.New(im, scales)
	log.Print("finished re-sampling levels: ", len(images.Levels))

	feats := make([]*rimg64.Multi, len(images.Levels))
	for i, x := range images.Levels {
		feats[i] = fn(rimg64.FromColor(x))
	}
	log.Print("finished computing feature transform: ", len(images.Levels))
	return &Pyramid{images, feats, rate}
}

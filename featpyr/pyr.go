package featpyr

import (
	"image"
	"log"

	"github.com/jvlmdr/go-cv/feat"
	"github.com/jvlmdr/go-cv/imgpyr"
	"github.com/jvlmdr/go-cv/rimg64"
)

// Level zero is the original image.
// Level i has size (width, height) * Scales.At(i).
type Pyramid struct {
	Images *imgpyr.Pyramid
	// Feature transform of each level in pyramid.
	Feats []*rimg64.Multi
	// Integer downsample rate of features.
	Rate int
	// Margin which was added around each image before computing features.
	Margin feat.Margin
}

// Constructs a feature pyramid.
// Extends each level by a margin before computing features.
func NewPad(images *imgpyr.Pyramid, phi feat.Transform, pad feat.Pad) *Pyramid {
	feats := make([]*rimg64.Multi, len(images.Levels))
	for i, im := range images.Levels {
		feats[i] = feat.ApplyPad(phi, im, pad)
	}
	log.Print("finished computing feature transform: ", len(images.Levels))
	return &Pyramid{images, feats, phi.Rate(), pad.Margin}
}

// Constructs a feature pyramid.
func New(images *imgpyr.Pyramid, phi feat.Transform) *Pyramid {
	return NewPad(images, phi, feat.NoPad())
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

// Converts a point in the feature pyramid to a point in the image.
func (pyr *Pyramid) ToImagePoint(pt imgpyr.Point) image.Point {
	return vec(pt.Pos.Mul(pyr.Rate)).Mul(1 / pyr.Scale(pt.Level)).Round()
}

// Converts a point in the feature pyramid to a rectangle in the image.
func (pyr *Pyramid) ToImageRect(pt imgpyr.Point, interior image.Rectangle) image.Rectangle {
	// Translate interior by position (scaled by rate) and subtract margin offset.
	rect := interior.Add(pt.Pos.Mul(pyr.Rate)).Sub(pyr.Margin.TopLeft())
	// Scale rectangle.
	return scaleRect(1/pyr.Scale(pt.Level), rect)
}

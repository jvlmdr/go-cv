package featpyr

import (
	"image"

	"github.com/jvlmdr/go-cv/feat"
	"github.com/jvlmdr/go-cv/imgpyr"
	"github.com/jvlmdr/go-cv/rimg64"
)

// Generator describes a feature pyramid which is generated level by level.
//
// The features are computed using Transform after extending the image using Pad.
type Generator struct {
	Image *imgpyr.Generator
	feat.Transform
	feat.Pad
}

func NewGenerator(im *imgpyr.Generator, phi feat.Transform, pad feat.Pad) *Generator {
	return &Generator{im, phi, pad}
}

// Level describes one level of a pyramid Generator.
type Level struct {
	Image *imgpyr.Level
	Feat  *rimg64.Multi
}

func (pyr *Generator) First() *Level {
	im := pyr.Image.First()
	if im == nil {
		return nil
	}
	x := feat.ApplyPad(pyr.Transform, im.Image, pyr.Pad)
	return &Level{im, x}
}

func (pyr *Generator) Next(curr *Level) *Level {
	im := pyr.Image.Next(curr.Image)
	if im == nil {
		return nil
	}
	x := feat.ApplyPad(pyr.Transform, im.Image, pyr.Pad)
	return &Level{im, x}
}

// ToImageRect converts a point in the feature pyramid to a rectangle in the image.
func (pyr *Generator) ToImageRect(level int, pt image.Point, interior image.Rectangle) image.Rectangle {
	// Translate interior by position (scaled by rate) and subtract margin offset.
	rate := pyr.Transform.Rate()
	offset := pyr.Pad.Margin.TopLeft()
	scale := pyr.Image.Scales[level]
	rect := interior.Add(pt.Mul(rate)).Sub(offset)
	return scaleRect(1/scale, rect)
}

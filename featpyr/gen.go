package featpyr

import (
	"image"
	"time"

	"github.com/jvlmdr/go-cv/feat"
	"github.com/jvlmdr/go-cv/imgpyr"
	"github.com/jvlmdr/go-cv/rimg64"
)

// Generator describes a feature pyramid which is generated level by level.
//
// The features are computed using Transform after extending the image using Pad.
type Generator struct {
	Image     *imgpyr.Generator
	Transform feat.Image
	feat.Pad
	// Cumulative time.
	DurResize time.Duration
	DurFeat   time.Duration
}

func NewGenerator(im *imgpyr.Generator, phi feat.Image, pad feat.Pad) *Generator {
	return &Generator{Image: im, Transform: phi, Pad: pad}
}

// Level describes one level of a pyramid Generator.
type Level struct {
	Image *imgpyr.Level
	Feat  *rimg64.Multi
}

func (pyr *Generator) First() (*Level, error) {
	t := time.Now()
	im := pyr.Image.First()
	pyr.DurResize = time.Since(t)
	if im == nil {
		return nil, nil
	}
	t = time.Now()
	x, err := feat.ApplyPad(pyr.Transform, im.Image, pyr.Pad)
	pyr.DurFeat = time.Since(t)
	if err != nil {
		return nil, err
	}
	return &Level{im, x}, nil
}

func (pyr *Generator) Next(curr *Level) (*Level, error) {
	t := time.Now()
	im := pyr.Image.Next(curr.Image)
	pyr.DurResize += time.Since(t)
	if im == nil {
		return nil, nil
	}
	t = time.Now()
	x, err := feat.ApplyPad(pyr.Transform, im.Image, pyr.Pad)
	pyr.DurFeat += time.Since(t)
	if err != nil {
		return nil, err
	}
	return &Level{im, x}, nil
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

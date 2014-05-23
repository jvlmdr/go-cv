package imgpyr

import (
	"image"
	"image/draw"
	"math"

	"github.com/nfnt/resize"
)

// Multi-resolution representation of an image.
type Pyramid struct {
	Levels []image.Image
	Scales GeoSeq
}

// A point in a pyramid.
type Point struct {
	Level int
	Pos   image.Point
}

// Returns set of scales for an image pyramid
// given minimum interesting image size.
// Step can be greater than or less than 1.
func Scales(im, tmpl image.Point, step float64) GeoSeq {
	if step <= 0 {
		panic("step must be positive")
	}
	// Rectify step to be between 0 and 1.
	if step > 1 {
		step = 1 / step
	}

	// Template dims as fraction of image dims.
	x := float64(tmpl.X) / float64(im.X)
	y := float64(tmpl.Y) / float64(im.Y)
	// Don't want to go smaller than either.
	return Sequence(1, step, math.Max(x, y))
}

// Default interpolation method used by New().
var DefaultInterp resize.InterpolationFunction = resize.Bicubic

// Creates a new pyramid.
//
// Interpolation method specified by DefaultInterp variable.
func New(img image.Image, scales GeoSeq) *Pyramid {
	return NewInterp(img, scales, DefaultInterp)
}

// Creates a new pyramid using specified interpolation.
func NewInterp(img image.Image, scales GeoSeq, interp resize.InterpolationFunction) *Pyramid {
	levels := make([]image.Image, scales.Len)
	levels[0] = clone(img)
	for i := 1; i < scales.Len; i++ {
		// Compute image dimensions at this level.
		scale := scales.At(i)
		width := round(float64(img.Bounds().Dx()) * scale)
		height := round(float64(img.Bounds().Dy()) * scale)
		levels[i] = resize.Resize(uint(width), uint(height), levels[i-1], interp)
	}
	return &Pyramid{levels, scales}
}

// Creates a copy of the image pyramid.
func (pyr *Pyramid) Clone() *Pyramid {
	levels := make([]image.Image, len(pyr.Levels))
	for i, level := range pyr.Levels {
		levels[i] = clone(level)
	}
	return &Pyramid{levels, pyr.Scales}
}

// Clones an image.Image, assuming it is an image.RGBA64.
func clone(img image.Image) image.Image {
	dst := image.NewRGBA64(img.Bounds())
	draw.Draw(dst, dst.Bounds(), img, image.ZP, draw.Src)
	return dst
}

func (pyr *Pyramid) Visualize() image.Image {
	var width, height int
	for _, level := range pyr.Levels {
		if level.Bounds().Dx() > width {
			width = level.Bounds().Dx()
		}
		height += level.Bounds().Dy()
	}

	img := image.NewRGBA64(image.Rect(0, 0, width, height))
	var p image.Point
	for _, level := range pyr.Levels {
		draw.Draw(img, level.Bounds().Add(p), level, image.ZP, draw.Src)
		p.Y += level.Bounds().Dy()
	}
	return img
}

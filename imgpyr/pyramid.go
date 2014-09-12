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
	Scales []float64
}

// A point in a pyramid.
type Point struct {
	Level int
	Pos   image.Point
}

// Returns set of scales for an image pyramid
// given minimum interesting image size.
// Step can be greater than or less than 1.
func Scales(im, tmpl image.Point, max, step float64) GeoSeq {
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
	return Sequence(max, step, math.Max(x, y))
}

// Default interpolation method used by New().
var DefaultInterp resize.InterpolationFunction = resize.Bicubic

// New creates a new pyramid.
//
// Interpolation method specified by DefaultInterp variable.
func New(im image.Image, scales []float64) *Pyramid {
	return NewInterp(im, scales, DefaultInterp)
}

func scaleSize(size image.Point, scale float64) image.Point {
	x := round(float64(size.X) * scale)
	y := round(float64(size.Y) * scale)
	return image.Pt(x, y)
}

// NewInterp creates a new pyramid using specified interpolation.
func NewInterp(im image.Image, scales []float64, interp resize.InterpolationFunction) *Pyramid {
	if len(scales) == 0 {
		return nil
	}
	levels := make([]image.Image, len(scales))
	orig := im.Bounds().Size()
	src := im
	for i := range scales {
		size := scaleSize(orig, scales[i])
		levels[i] = resizeIfNec(size, src, interp)
		if scales[i] < 1 {
			// If the current level is smaller than the original
			// then use it as the source for the next level.
			src = levels[i]
		}
	}
	return &Pyramid{levels, scales}
}

// Clones the image if resize is not necessary.
func resizeIfNec(size image.Point, im image.Image, interp resize.InterpolationFunction) image.Image {
	if size.Eq(im.Bounds().Size()) {
		return clone(im)
	}
	return resize.Resize(uint(size.X), uint(size.Y), im, interp)
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
func clone(im image.Image) image.Image {
	dst := image.NewRGBA64(im.Bounds())
	draw.Draw(dst, dst.Bounds(), im, image.ZP, draw.Src)
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

	im := image.NewRGBA64(image.Rect(0, 0, width, height))
	var p image.Point
	for _, level := range pyr.Levels {
		draw.Draw(im, level.Bounds().Add(p), level, image.ZP, draw.Src)
		p.Y += level.Bounds().Dy()
	}
	return im
}

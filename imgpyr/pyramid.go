package imgpyr

import (
	"github.com/nfnt/resize"

	"fmt"
	"image"
	"image/draw"
	"math"
)

// Default interpolation method used by New().
var DefaultInterp resize.InterpolationFunction = resize.Bicubic

type Pyramid struct {
	Levels []image.Image
	Scales GeoSeq
}

// Describes a finite geometric sequence.
type GeoSeq struct {
	Start float64
	Step  float64
	Len   int
}

// Returns the i-th value of the progression.
func (seq GeoSeq) At(i int) float64 {
	if i < 0 || i >= seq.Len {
		panic(fmt.Sprintf("out of range: %d", i))
	}
	return seq.Start * math.Pow(seq.Step, float64(i))
}

// Returns the (floating point) index of the x in the progression.
func (seq GeoSeq) Inv(x float64) float64 {
	return math.Log(x / seq.Start) / math.Log(seq.Step)
}

// Creates a new pyramid.
//
// Interpolation method specified by DefaultInterp variable.
func New(img image.Image, scales GeoSeq) *Pyramid {
	return NewInterp(img, scales, DefaultInterp)
}

// Creates a new pyramid using specified interpolation.
func NewInterp(img image.Image, scales GeoSeq, interp resize.InterpolationFunction) *Pyramid {
	levels := make([]image.Image, scales.Len)
	for i := 0; i < scales.Len; i++ {
		// Compute image dimensions at this level.
		scale := scales.At(i)
		width := round(float64(img.Bounds().Dx()) * scale)
		height := round(float64(img.Bounds().Dy()) * scale)
		levels[i] = resize.Resize(uint(width), uint(height), img, interp)
	}
	return &Pyramid{levels, scales}
}

// Returns the nearest integer.
func round(x float64) int {
	return int(math.Floor(x+0.5))
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

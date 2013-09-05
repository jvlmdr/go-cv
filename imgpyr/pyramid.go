package imgpyr

import (
	"github.com/nfnt/resize"

	"fmt"
	"image"
	"image/draw"
	"math"
)

// Multi-resolution representation of an image.
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

// Generates a sequence from start to the last element <= lim.
//
// If step > 1, then lim must be greater than start.
// If step < 1, then lim must be less than start.
func Sequence(start, step, lim float64) GeoSeq {
	n := math.Log(lim/start) / math.Log(step)
	return GeoSeq{start, step, int(math.Floor(n)) + 1}
}

// Generates a sequence from first to last containing n elements.
func LogRange(first, last float64, n int) GeoSeq {
	step := math.Exp(math.Log(last/first) / float64(n-1))
	return GeoSeq{first, step, n}
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
	return math.Log(x/seq.Start) / math.Log(seq.Step)
}

// Returns a reversed sequence.
func (seq GeoSeq) Reverse() GeoSeq {
	return GeoSeq{seq.At(seq.Len - 1), 1 / seq.Step, seq.Len}
}

// Default interpolation method used by New().
var DefaultInterp resize.InterpolationFunction = resize.Bicubic

// Creates a new pyramid.
//
// Interpolation method specified by DefaultInterp variable.
func New(img image.Image, scales GeoSeq) *Pyramid {
	// Ensure that images go from smallest to largest.
	if scales.Step < 1 {
		scales = scales.Reverse()
	}
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
	return int(math.Floor(x + 0.5))
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

package rimg64

import (
	"image"
	"image/color"
	"math"

	"github.com/gonum/floats"
)

// Image describes a real-valued image.
type Image struct {
	// Element (x, y) at index x*Height + y.
	Elems  []float64
	Width  int
	Height int
}

// New allocates an image of zeros.
func New(width, height int) *Image {
	pixels := make([]float64, width*height)
	return &Image{pixels, width, height}
}

// Size returns the dimensions of the image.
func (f *Image) Size() image.Point {
	return image.Pt(f.Width, f.Height)
}

// At retrieves an element of the image.
func (f *Image) At(x, y int) float64 {
	return f.Elems[x*f.Height+y]
}

// Set modifies an element of the image.
func (f *Image) Set(x, y int, v float64) {
	f.Elems[f.index(x, y)] = v
}

func (f *Image) index(x, y int) int {
	return x*f.Height + y
}

// Clone creates a copy of the image.
func (f *Image) Clone() *Image {
	g := New(f.Width, f.Height)
	copy(g.Elems, f.Elems)
	return g
}

// SubImage clones part of an image.
func (f *Image) SubImage(r image.Rectangle) *Image {
	g := New(r.Dx(), r.Dy())
	for i := 0; i < g.Width; i++ {
		for j := 0; j < g.Height; j++ {
			g.Set(i, j, f.At(r.Min.X+i, r.Min.Y+j))
		}
	}
	return g
}

// ToGray converts the image to an 8-bit integer gray image.
// Maps [0, 1] to [0, math.MaxUint8].
// Caps to [0, math.MaxUint8].
func ToGray(f *Image) *image.Gray {
	g := image.NewGray(image.Rect(0, 0, f.Width, f.Height))
	for x := 0; x < f.Width; x++ {
		for y := 0; y < f.Height; y++ {
			v := round(math.MaxUint8 * f.At(x, y))
			v = min(math.MaxUint8, max(0, v))
			g.SetGray(x, y, color.Gray{uint8(v)})
		}
	}
	return g
}

// ToGray16 converts the image to a 16-bit integer gray image.
// Maps [0, 1] to [0, math.MaxUint16].
// Caps to [0, math.MaxUint16].
func ToGray16(f *Image) *image.Gray16 {
	g := image.NewGray16(image.Rect(0, 0, f.Width, f.Height))
	for x := 0; x < f.Width; x++ {
		for y := 0; y < f.Height; y++ {
			v := round(math.MaxUint16 * f.At(x, y))
			v = min(math.MaxUint16, max(0, v))
			g.SetGray16(x, y, color.Gray16{uint16(v)})
		}
	}
	return g
}

// FromGray converts the image to a real-valued gray image with values in [0, 1].
func FromGray(g image.Image) *Image {
	size := g.Bounds().Size()
	off := g.Bounds().Min
	f := New(size.X, size.Y)
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			v := color.Gray16Model.Convert(g.At(x+off.X, y+off.Y))
			c, _, _, _ := v.RGBA()
			f.Set(x, y, float64(c)/float64(0xFFFF))
		}
	}
	return f
}

// Plus computes the sum of two images.
// Does not modify either input.
func (f *Image) Plus(g *Image) *Image {
	dst := New(f.Width, f.Height)
	floats.AddTo(dst.Elems, f.Elems, g.Elems)
	return dst
}

// Minus computes the difference between two images.
// Does not modify either input.
func (f *Image) Minus(g *Image) *Image {
	dst := New(f.Width, f.Height)
	floats.SubTo(dst.Elems, f.Elems, g.Elems)
	return dst
}

// Scale computes the product of an image with a scalar.
// Does not modify the input.
func (f *Image) Scale(alpha float64) *Image {
	dst := New(f.Width, f.Height)
	floats.AddScaled(dst.Elems, alpha, f.Elems)
	return dst
}

// FromRows creates a new image from a list of rows.
// Panics if rows are varying length.
// Returns nil if no rows are provided.
// Returns empty image if all rows are empty.
func FromRows(rows [][]float64) *Image {
	height := len(rows)
	if height == 0 {
		return nil
	}
	width := len(rows[0])
	f := New(width, height)
	for j, row := range rows {
		if len(row) != width {
			panic("rows are not same length")
		}
		for i, x := range row {
			f.Set(i, j, x)
		}
	}
	return f
}

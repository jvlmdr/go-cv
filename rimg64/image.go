package rimg64

import (
	"image"
	"image/color"
)

// Describes a real-valued image.
type Image struct {
	// Element (x, y) at index x*Height + y.
	Elems  []float64
	Width  int
	Height int
}

// Allocates an image of zeros.
func New(width, height int) *Image {
	pixels := make([]float64, width*height)
	return &Image{pixels, width, height}
}

func (f *Image) Size() image.Point {
	return image.Pt(f.Width, f.Height)
}

func (f *Image) At(x, y int) float64 {
	return f.Elems[x*f.Height+y]
}

func (f *Image) Set(x, y int, v float64) {
	f.Elems[f.index(x, y)] = v
}

func (f *Image) index(x, y int) int {
	return x*f.Height + y
}

// Creates a copy of the image.
func (f *Image) Clone() *Image {
	g := New(f.Width, f.Height)
	copy(g.Elems, f.Elems)
	return g
}

// Converts to an 8-bit integer gray image.
// Maps [0, 1] to [0, 255].
// Caps to [0, 255].
func ToGray(f *Image) *image.Gray {
	g := image.NewGray(image.Rect(0, 0, f.Width, f.Height))
	for x := 0; x < f.Width; x++ {
		for y := 0; y < f.Height; y++ {
			v := round(255 * f.At(x, y))
			v = min(255, max(0, v))
			g.SetGray(x, y, color.Gray{uint8(v)})
		}
	}
	return g
}

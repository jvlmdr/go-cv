package rimg64

import (
	"image"
	"image/color"

	"github.com/gonum/floats"
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

// Clones an image from part of a larger image.
func (f *Image) SubImage(r image.Rectangle) *Image {
	g := New(r.Dx(), r.Dy())
	for i := 0; i < g.Width; i++ {
		for j := 0; j < g.Height; j++ {
			g.Set(i, j, f.At(r.Min.X+i, r.Min.Y+j))
		}
	}
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

// Converts a color image to a 3-channel vector-valued image.
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

// Returns the sum of two images.
// Does not modify either input.
func (f *Image) Plus(g *Image) *Image {
	dst := New(f.Width, f.Height)
	floats.Add(dst.Elems, f.Elems, g.Elems)
	return dst
}

// Returns the difference of two images.
// Does not modify either input.
func (f *Image) Minus(g *Image) *Image {
	dst := New(f.Width, f.Height)
	floats.SubTo(dst.Elems, f.Elems, g.Elems)
	return dst
}

// Returns a scaled copy of an image.
func (f *Image) Scale(alpha float64) *Image {
	dst := New(f.Width, f.Height)
	floats.AddScaled(dst.Elems, alpha, f.Elems)
	return dst
}

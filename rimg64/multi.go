package rimg64

import (
	"fmt"
	"image"

	"github.com/gonum/floats"
)

// Multi describes a multi-channel (real vector valued) image.
type Multi struct {
	// Element (x, y, d) at index (x*f.Height + y)*f.Channels + d.
	Elems    []float64
	Width    int
	Height   int
	Channels int
}

// NewMulti allocates a multi-channel image of zeros.
func NewMulti(width, height, channels int) *Multi {
	pixels := make([]float64, width*height*channels)
	return &Multi{pixels, width, height, channels}
}

// Size returns the width and height of the image.
func (f *Multi) Size() image.Point {
	return image.Pt(f.Width, f.Height)
}

// At retrieves an element of the image.
func (f *Multi) At(x, y, d int) float64 {
	return f.Elems[f.index(x, y, d)]
}

// Set modifies an element of the image.
func (f *Multi) Set(x, y, d int, v float64) {
	f.Elems[f.index(x, y, d)] = v
}

func (f *Multi) index(x, y, d int) int {
	return (x*f.Height+y)*f.Channels + d
}

// Clone creates a copy of the image.
func (f *Multi) Clone() *Multi {
	g := NewMulti(f.Width, f.Height, f.Channels)
	copy(g.Elems, f.Elems)
	return g
}

// Pixel retrieves all channels at a given position in the image.
// The returned slice does not reference the elements of the original image.
func (f *Multi) Pixel(x, y int) []float64 {
	v := make([]float64, f.Channels)
	for d := range v {
		v[d] = f.At(x, y, d)
	}
	return v
}

// SetPixel modifies all channels at a given position in the image.
// Panics if the number of elements does not match the number of channels.
func (f *Multi) SetPixel(x, y int, v []float64) {
	if len(v) != f.Channels {
		panic(fmt.Sprintf("different number of channels: image %d, vector %d", f.Channels, len(v)))
	}
	for d := 0; d < f.Channels; d++ {
		f.Set(x, y, d, v[d])
	}
}

// Channel retrieves all elements in one channel as a real-valued image.
// The returned image does not reference the elements of the original image.
func (f *Multi) Channel(d int) *Image {
	g := New(f.Width, f.Height)
	for x := 0; x < f.Width; x++ {
		for y := 0; y < f.Height; y++ {
			g.Set(x, y, f.At(x, y, d))
		}
	}
	return g
}

// SetChannel modifies all elements in one channel.
// Panics if the sizes do not match.
func (f *Multi) SetChannel(d int, g *Image) {
	if !f.Size().Eq(g.Size()) {
		panic(fmt.Sprintf("different size: image %v, channel %v", f.Size(), g.Size()))
	}
	for x := 0; x < f.Width; x++ {
		for y := 0; y < f.Height; y++ {
			f.Set(x, y, d, g.At(x, y))
		}
	}
}

// SubImage clones part of an image.
func (f *Multi) SubImage(r image.Rectangle) *Multi {
	g := NewMulti(r.Dx(), r.Dy(), f.Channels)
	for i := 0; i < g.Width; i++ {
		for j := 0; j < g.Height; j++ {
			for k := 0; k < g.Channels; k++ {
				g.Set(i, j, k, f.At(r.Min.X+i, r.Min.Y+j, k))
			}
		}
	}
	return g
}

// FromColor converts a color image into a 3-channel real-valued image.
func FromColor(g image.Image) *Multi {
	size := g.Bounds().Size()
	off := g.Bounds().Min
	f := NewMulti(size.X, size.Y, 3)
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			var c [3]uint32
			c[0], c[1], c[2], _ = g.At(x+off.X, y+off.Y).RGBA()
			for d := 0; d < 3; d++ {
				f.Set(x, y, d, float64(c[d])/float64(0xFFFF))
			}
		}
	}
	return f
}

// Plus computes the sum of two images.
// Does not modify either input.
func (f *Multi) Plus(g *Multi) *Multi {
	dst := NewMulti(f.Width, f.Height, f.Channels)
	floats.Add(dst.Elems, f.Elems, g.Elems)
	return dst
}

// Minus computes the difference between two images.
// Does not modify either input.
func (f *Multi) Minus(g *Multi) *Multi {
	dst := NewMulti(f.Width, f.Height, f.Channels)
	floats.SubTo(dst.Elems, f.Elems, g.Elems)
	return dst
}

// Scale computes the product of an image with a scalar.
// Does not modify the input.
func (f *Multi) Scale(alpha float64) *Multi {
	dst := NewMulti(f.Width, f.Height, f.Channels)
	floats.AddScaled(dst.Elems, alpha, f.Elems)
	return dst
}

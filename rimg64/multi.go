package rimg64

import "image"

// Describe a real-vector-valued image.
type Multi struct {
	// Element (x, y, d) at index (x*f.Height + y)*f.Channels + d.
	Elems    []float64
	Width    int
	Height   int
	Channels int
}

// Allocates an image of zeros.
func NewMulti(width, height, channels int) *Multi {
	pixels := make([]float64, width*height*channels)
	return &Multi{pixels, width, height, channels}
}

func (f *Multi) Size() image.Point {
	return image.Pt(f.Width, f.Height)
}

func (f *Multi) At(x, y, d int) float64 {
	return f.Elems[f.index(x, y, d)]
}

func (f *Multi) Set(x, y, d int, v float64) {
	f.Elems[f.index(x, y, d)] = v
}

func (f *Multi) index(x, y, d int) int {
	return (x*f.Height+y)*f.Channels + d
}

// Creates a copy of the image.
func (f *Multi) Clone() *Multi {
	g := NewMulti(f.Width, f.Height, f.Channels)
	copy(g.Elems, f.Elems)
	return g
}

// Clones the value of all channels at a point.
func (f *Multi) Pixel(x, y int) []float64 {
	v := make([]float64, f.Channels)
	for d := range v {
		v[d] = f.At(x, y, d)
	}
	return v
}

// Overwrites the value of all channels at a point.
func (f *Multi) SetPixel(x, y int, v []float64) {
	for d := 0; d < f.Channels; d++ {
		f.Set(x, y, d, v[d])
	}
}

// Clones one channel of a vector image.
func (f *Multi) Channel(d int) *Image {
	g := New(f.Width, f.Height)
	for x := 0; x < f.Width; x++ {
		for y := 0; y < f.Height; y++ {
			g.Set(x, y, f.At(x, y, d))
		}
	}
	return g
}

// Overwrites one channel with a scalar image.
func (f *Multi) SetChannel(d int, g *Image) {
	for x := 0; x < f.Width; x++ {
		for y := 0; y < f.Height; y++ {
			f.Set(x, y, d, g.At(x, y))
		}
	}
}

// Clones an image from part of a larger image.
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

// Converts a color image to a 3-channel vector-valued image.
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

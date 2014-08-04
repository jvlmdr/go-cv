package rimg64

import "image"

type Table Image

// CumSum computes a summed area table or integral image.
func CumSum(f *Image) *Table {
	s := New(f.Width, f.Height)
	for i := 0; i < f.Width; i++ {
		for j := 0; j < f.Height; j++ {
			t := f.At(i, j)
			if i > 0 {
				t += s.At(i-1, j)
			}
			if j > 0 {
				t += s.At(i, j-1)
			}
			if i > 0 && j > 0 {
				t -= s.At(i-1, j-1)
			}
			s.Set(i, j, t)
		}
	}
	return (*Table)(s)
}

// Rect returns the sum of elements in the region
// 	r.Min.X <= x < r.Max.X
// 	r.Min.Y <= y < r.Max.Y
func (t *Table) Rect(r image.Rectangle) float64 {
	s := (*Image)(t)
	bnds := image.Rect(0, 0, s.Width, s.Height)
	if !r.In(bnds) {
		panic("out of bounds")
	}
	if r.Dx()*r.Dy() == 0 {
		return 0
	}
	area := s.At(r.Max.X-1, r.Max.Y-1)
	if r.Min.X > 0 {
		area -= s.At(r.Min.X-1, r.Max.Y-1)
	}
	if r.Min.Y > 0 {
		area -= s.At(r.Max.X-1, r.Min.Y-1)
	}
	if r.Min.X > 0 && r.Min.Y > 0 {
		area += s.At(r.Min.X-1, r.Min.Y-1)
	}
	return area
}

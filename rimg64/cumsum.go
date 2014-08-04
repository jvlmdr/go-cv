package rimg64

// CumSum computes a summed area table or integral image.
func CumSum(f *Image) *Image {
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
	return s
}

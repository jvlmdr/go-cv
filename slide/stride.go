package slide

import (
	"github.com/gonum/floats"
	"github.com/jvlmdr/go-cv/rimg64"
)

// Decimate takes every k-th sample starting at (0, 0).
func Decimate(f *rimg64.Image, k int) *rimg64.Image {
	g := rimg64.New(ceilDiv(f.Width, k), ceilDiv(f.Height, k))
	for i := 0; i < g.Width; i++ {
		for j := 0; j < g.Height; j++ {
			g.Set(i, j, f.At(k*i, k*j))
		}
	}
	return g
}

// CorrStride is a more efficient way to compute Decimate(Corr(f, g), k).
func CorrStride(f, g *rimg64.Image, k int) *rimg64.Image {
	// Convolution and downsampling can be expressed as
	// the sum of convolutions of downsampled signals.

	// Compute size of downsampled output.
	dst := ValidSizeStride(f.Size(), g.Size(), k)
	// Grid of convolutions to perform is m x n.
	m, n := min(k, g.Width), min(k, g.Height)

	// Divide g into grid of downsampled signals.
	gs := make([][]*rimg64.Image, m)
	for i := 0; i < m; i++ {
		gs[i] = make([]*rimg64.Image, n)
		for j := 0; j < n; j++ {
			// kt + i <= l - 1
			// kt <= l - 1 - i
			// t <= floor((l-1-i) / k)
			w := (g.Width-1-i)/k + 1
			h := (g.Height-1-j)/k + 1
			gs[i][j] = rimg64.New(w, h)
		}
	}
	for u := 0; u < g.Width; u++ {
		for v := 0; v < g.Height; v++ {
			gs[u%k][v%k].Set(u/k, v/k, g.At(u, v))
		}
	}

	// Divide f into grid of downsampled signals.
	fs := make([][]*rimg64.Image, len(gs))
	for i := 0; i < m; i++ {
		fs[i] = make([]*rimg64.Image, len(gs[i]))
		for j := 0; j < n; j++ {
			// The size of each f is such that convolution with
			// the corresponding g yields the output size.
			w := dst.X + gs[i][j].Width - 1
			h := dst.Y + gs[i][j].Height - 1
			fs[i][j] = rimg64.New(w, h)
			// Populate image.
			for u := 0; u < w; u++ {
				for v := 0; v < h; v++ {
					fs[i][j].Set(u, v, f.At(k*u+i, k*v+j))
				}
			}
		}
	}

	h := rimg64.New(dst.X, dst.Y)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			hij := Corr(fs[i][j], gs[i][j])
			floats.Add(h.Elems, hij.Elems)
		}
	}
	return h
}

// CorrMultiStride is a more efficient way to compute Decimate(CorrMulti(f, g), k).
func CorrMultiStride(f, g *rimg64.Multi, k int) *rimg64.Image {
	if err := errIfChannelsNotEq(f, g); err != nil {
		panic(err)
	}

	// Compute size of downsampled output.
	dst := ValidSizeStride(f.Size(), g.Size(), k)
	// Grid of convolutions to perform is m x n.
	m, n := min(k, g.Width), min(k, g.Height)

	// Divide g into grid of downsampled signals.
	gs := make([][]*rimg64.Multi, m)
	for i := 0; i < m; i++ {
		gs[i] = make([]*rimg64.Multi, n)
		for j := 0; j < n; j++ {
			// kt + i <= l - 1
			// kt <= l - 1 - i
			// t <= floor((l-1-i) / k)
			w := (g.Width-1-i)/k + 1
			h := (g.Height-1-j)/k + 1
			gs[i][j] = rimg64.NewMulti(w, h, g.Channels)
		}
	}
	for u := 0; u < g.Width; u++ {
		for v := 0; v < g.Height; v++ {
			for p := 0; p < g.Channels; p++ {
				gs[u%k][v%k].Set(u/k, v/k, p, g.At(u, v, p))
			}
		}
	}

	// Divide f into grid of downsampled signals.
	fs := make([][]*rimg64.Multi, len(gs))
	for i := 0; i < m; i++ {
		fs[i] = make([]*rimg64.Multi, len(gs[i]))
		for j := 0; j < n; j++ {
			// The size of each f is such that convolution with
			// the corresponding g yields the output size.
			w := dst.X + gs[i][j].Width - 1
			h := dst.Y + gs[i][j].Height - 1
			fs[i][j] = rimg64.NewMulti(w, h, f.Channels)
			// Populate image.
			for u := 0; u < w; u++ {
				for v := 0; v < h; v++ {
					for p := 0; p < f.Channels; p++ {
						fs[i][j].Set(u, v, p, f.At(k*u+i, k*v+j, p))
					}
				}
			}
		}
	}

	h := rimg64.New(dst.X, dst.Y)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			hij := CorrMulti(fs[i][j], gs[i][j])
			floats.Add(h.Elems, hij.Elems)
		}
	}
	return h
}

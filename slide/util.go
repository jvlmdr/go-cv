package slide

import (
	"fmt"
	"image"
	"math/cmplx"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-fftw/fftw"
)

func copyImageTo(x *fftw.Array2, f *rimg64.Image) {
	w, h := x.Dims()
	for u := 0; u < w; u++ {
		for v := 0; v < h; v++ {
			if u < f.Width && v < f.Height {
				x.Set(u, v, complex(f.At(u, v), 0))
			} else {
				x.Set(u, v, 0)
			}
		}
	}
}

// Assumes that f is no smaller than x.
// Pads with zeros.
func copyChannelTo(x *fftw.Array2, f *rimg64.Multi, p int) {
	w, h := x.Dims()
	for u := 0; u < w; u++ {
		for v := 0; v < h; v++ {
			if u < f.Width && v < f.Height {
				x.Set(u, v, complex(f.At(u, v, p), 0))
			} else {
				x.Set(u, v, 0)
			}
		}
	}
}

// Assumes that f is no smaller than x.
func copyRealTo(f *rimg64.Image, x *fftw.Array2) {
	for u := 0; u < f.Width; u++ {
		for v := 0; v < f.Height; v++ {
			f.Set(u, v, real(x.At(u, v)))
		}
	}
}

// Assumes that f is no smaller than x.
func copyRealToChannel(f *rimg64.Multi, p int, x *fftw.Array2) {
	for u := 0; u < f.Width; u++ {
		for v := 0; v < f.Height; v++ {
			f.Set(u, v, p, real(x.At(u, v)))
		}
	}
}

// dst[i, j] = src[i*stride + offset.X, j*stride + offset.Y],
// or zero if this is outside the boundary.
func copyStrideTo(dst *fftw.Array2, src *rimg64.Image, stride int, offset image.Point) {
	m, n := dst.Dims()
	bnds := image.Rect(0, 0, src.Width, src.Height)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			p := image.Pt(i, j).Mul(stride).Add(offset)
			var val complex128
			if p.In(bnds) {
				val = complex(src.At(p.X, p.Y), 0)
			}
			dst.Set(i, j, val)
		}
	}
}

// dst[i, j] = src[i*stride + offset.X, j*stride + offset.Y],
// or zero if this is outside the boundary.
func copyChannelStrideTo(dst *fftw.Array2, src *rimg64.Multi, channel, stride int, offset image.Point) {
	m, n := dst.Dims()
	bnds := image.Rect(0, 0, src.Width, src.Height)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			p := image.Pt(i, j).Mul(stride).Add(offset)
			var val complex128
			if p.In(bnds) {
				val = complex(src.At(p.X, p.Y, channel), 0)
			}
			dst.Set(i, j, val)
		}
	}
}

// ci <- ci + k conj(ai) bi
func addScaleMul(c *fftw.Array2, k complex128, a, b *fftw.Array2) {
	m, n := a.Dims()
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			ax := cmplx.Conj(a.At(i, j))
			bx := b.At(i, j)
			cx := c.At(i, j)
			c.Set(i, j, cx+k*ax*bx)
		}
	}
}

// ci <- k conj(ai) bi
func scaleMul(c *fftw.Array2, k complex128, a, b *fftw.Array2) {
	m, n := a.Dims()
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			ax := cmplx.Conj(a.At(i, j))
			bx := b.At(i, j)
			c.Set(i, j, k*ax*bx)
		}
	}
}

// ai <- k ai
func scale(k complex128, c *fftw.Array2) {
	m, n := c.Dims()
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			c.Set(i, j, k*c.At(i, j))
		}
	}
}

// ci <- ci + conj(ai) bi
func addMul(c *fftw.Array2, a, b *fftw.Array2) {
	m, n := a.Dims()
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			ax := cmplx.Conj(a.At(i, j))
			bx := b.At(i, j)
			cx := c.At(i, j)
			c.Set(i, j, cx+ax*bx)
		}
	}
}

func zero(x *fftw.Array2) {
	m, n := x.Dims()
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			x.Set(i, j, 0)
		}
	}
}

func errIfChannelsNotEq(f, g *rimg64.Multi) error {
	if f.Channels != g.Channels {
		return fmt.Errorf("different number of channels: %d, %d", f.Channels, g.Channels)
	}
	return nil
}

package slide

import (
	"image"
	"math/cmplx"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-fftw/fftw"
	"github.com/jvlmdr/lin-go/blas"
)

// MultiBank describes a collection of multi-channel filters.
// All filters must have the same dimension.
type MultiBank struct {
	Width    int
	Height   int
	Channels int
	Filters  []*rimg64.Multi
}

func (bank *MultiBank) Size() image.Point {
	return image.Pt(bank.Width, bank.Height)
}

func fftOfChannel(fhat *fftw.Array2, f *rimg64.Multi, p int) {
	plan := fftw.NewPlan2(fhat, fhat, fftw.Forward, fftw.Estimate)
	defer plan.Destroy()
	copyChannelTo(fhat, f, p)
	plan.Execute()
}

// ci <- ci + k ai* bi
func mulAddTo(c *fftw.Array2, k complex128, a, b *fftw.Array2) {
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

func zero(x *fftw.Array2) {
	m, n := x.Dims()
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			x.Set(i, j, 0)
		}
	}
}

// h_p[u, v] = sum_q (f_q corr g_pq)[u, v]
func CorrMultiBankNaive(f *rimg64.Multi, g *MultiBank) (*rimg64.Multi, error) {
	out := ValidSize(f.Size(), g.Size())
	if out.X <= 0 || out.Y <= 0 {
		return nil, nil
	}
	h := rimg64.NewMulti(out.X, out.Y, len(g.Filters))
	for u := 0; u < h.Width; u++ {
		for v := 0; v < h.Height; v++ {
			for p := 0; p < h.Channels; p++ {
				for i := 0; i < g.Width; i++ {
					for j := 0; j < g.Height; j++ {
						for q := 0; q < g.Channels; q++ {
							val := f.At(i+u, j+v, q) * g.Filters[p].At(i, j, q)
							h.Set(u, v, p, h.At(u, v, p)+val)
						}
					}
				}
			}
		}
	}
	return h, nil
}

// h_p[u, v] = sum_q (f_q corr g_pq)[u, v]
func CorrMultiBankFFT(f *rimg64.Multi, g *MultiBank) (*rimg64.Multi, error) {
	out := ValidSize(f.Size(), g.Size())
	if out.X <= 0 || out.Y <= 0 {
		return nil, nil
	}
	// Determine optimal size for FFT.
	work, _ := FFT2Size(f.Size())
	// Cache FFT of each channel of image.
	fhat := make([]*fftw.Array2, f.Channels)
	for i := range fhat {
		fhat[i] = fftw.NewArray2(work.X, work.Y)
		fftOfChannel(fhat[i], f, i)
	}

	curr := fftw.NewArray2(work.X, work.Y)
	fwd := fftw.NewPlan2(curr, curr, fftw.Forward, fftw.Estimate)
	defer fwd.Destroy()
	sum := fftw.NewArray2(work.X, work.Y)
	bwd := fftw.NewPlan2(sum, sum, fftw.Backward, fftw.Estimate)
	defer bwd.Destroy()

	h := rimg64.NewMulti(out.X, out.Y, len(g.Filters))
	alpha := complex(1/float64(work.X*work.Y), 0)
	// For each output channel.
	for p, gp := range g.Filters {
		zero(sum)
		// For each input channel.
		for q := 0; q < f.Channels; q++ {
			// Take FFT of this input channel.
			copyChannelTo(curr, gp, q)
			fwd.Execute()
			// h_p[x] = (F_p corr G_qp)[x]
			// H_p[x] = F_p[x]* G_qp[x]
			mulAddTo(sum, alpha, curr, fhat[q])
		}
		bwd.Execute()
		copyRealToChannel(h, p, sum)
	}
	return h, nil
}

func CorrMultiBankBLAS(f *rimg64.Multi, g *MultiBank) (*rimg64.Multi, error) {
	out := ValidSize(f.Size(), g.Size())
	if out.X <= 0 || out.Y <= 0 {
		return nil, nil
	}
	// Express as dense matrix multiplication.
	//   h_p[u, v] = sum_q (f_q corr g_pq)[u, v]
	//   h = A(f) X(g)
	// where A is (M-m+1)(N-n+1)k by mnk.

	h := rimg64.NewMulti(out.X, out.Y, len(g.Filters))
	M, N, K := h.Width, h.Height, h.Channels
	m, n, k := g.Width, g.Height, g.Channels
	a := blas.NewMat(M*N, m*n*k)
	{
		var r int
		for u := 0; u < h.Width; u++ {
			for v := 0; v < h.Height; v++ {
				var s int
				for i := 0; i < g.Width; i++ {
					for j := 0; j < g.Height; j++ {
						for q := 0; q < g.Channels; q++ {
							a.Set(r, s, f.At(i+u, j+v, q))
							s++
						}
					}
				}
				r++
			}
		}
	}
	x := blas.NewMat(m*n*k, K)
	{
		var r int
		for i := 0; i < g.Width; i++ {
			for j := 0; j < g.Height; j++ {
				for q := 0; q < g.Channels; q++ {
					for p := 0; p < h.Channels; p++ {
						x.Set(r, p, g.Filters[p].At(i, j, q))
					}
					r++
				}
			}
		}
	}
	y := blas.MatMul(1, a, x)
	{
		var r int
		for u := 0; u < h.Width; u++ {
			for v := 0; v < h.Height; v++ {
				for p := 0; p < h.Channels; p++ {
					h.Set(u, v, p, y.At(r, p))
				}
				r++
			}
		}
	}
	return h, nil
}

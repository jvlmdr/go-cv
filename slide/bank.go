package slide

import (
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-fftw/fftw"
	"github.com/jvlmdr/lin-go/blas"
)

// Bank describes a collection of filters.
// All filters must have the same spatial dimension.
type Bank struct {
	Width   int
	Height  int
	Filters []*rimg64.Image
}

// Size gives the spatial dimension of all filters in the bank.
func (bank *Bank) Size() image.Point {
	return image.Pt(bank.Width, bank.Height)
}

// CorrBankNaive computes the correlation of an image with a bank of filters.
// 	h_p[u, v] = (f corr g_p)[u, v]
func CorrBankNaive(f *rimg64.Image, g *Bank) (*rimg64.Multi, error) {
	out := ValidSize(f.Size(), g.Size())
	if out.X <= 0 || out.Y <= 0 {
		return nil, nil
	}
	h := rimg64.NewMulti(out.X, out.Y, len(g.Filters))
	for u := 0; u < h.Width; u++ {
		for v := 0; v < h.Height; v++ {
			for p := 0; p < h.Channels; p++ {
				var total float64
				for i := 0; i < g.Width; i++ {
					for j := 0; j < g.Height; j++ {
						total += f.At(i+u, j+v) * g.Filters[p].At(i, j)
					}
				}
				h.Set(u, v, p, total)
			}
		}
	}
	return h, nil
}

// CorrBankFFT computes the correlation of an image with a bank of filters.
// 	h_p[u, v] = (f corr g_p)[u, v]
func CorrBankFFT(f *rimg64.Image, g *Bank) (*rimg64.Multi, error) {
	out := ValidSize(f.Size(), g.Size())
	if out.X <= 0 || out.Y <= 0 {
		return nil, nil
	}
	// Determine optimal size for FFT.
	work, _ := FFT2Size(f.Size())
	// Re-use FFT of image.
	fhat := fftw.NewArray2(work.X, work.Y)
	copyImageTo(fhat, f)
	fftw.FFT2To(fhat, fhat)
	// Transform of each filter.
	curr := fftw.NewArray2(work.X, work.Y)
	fwd := fftw.NewPlan2(curr, curr, fftw.Forward, fftw.Estimate)
	defer fwd.Destroy()
	bwd := fftw.NewPlan2(curr, curr, fftw.Backward, fftw.Estimate)
	defer bwd.Destroy()

	h := rimg64.NewMulti(out.X, out.Y, len(g.Filters))
	alpha := complex(1/float64(work.X*work.Y), 0)
	// For each output channel.
	for p, gp := range g.Filters {
		// Take FFT.
		copyImageTo(curr, gp)
		fwd.Execute()
		// h_p[x] = (G_p corr F)[x]
		// H_p[x] = conj(G_p[x]) F[x]
		scaleMul(curr, alpha, curr, fhat)
		bwd.Execute()
		copyRealToChannel(h, p, curr)
	}
	return h, nil
}

// CorrBankBLAS computes the correlation of an image with a bank of filters.
// 	h_p[u, v] = (f corr g_p)[u, v]
func CorrBankBLAS(f *rimg64.Image, g *Bank) (*rimg64.Multi, error) {
	out := ValidSize(f.Size(), g.Size())
	if out.X <= 0 || out.Y <= 0 {
		return nil, nil
	}
	// Express as dense matrix multiplication.
	//   h_p[u, v] = (f corr g_q)[u, v]
	//   Y(h) = A(f) X(g)
	// If the number of output channels is k, then
	//   A is (M-m+1)(N-n+1) x mn and
	//   X is mn x k, so that
	//   Y is (M-m+1)(N-n+1) x k.

	h := rimg64.NewMulti(out.X, out.Y, len(g.Filters))
	m, n, k := g.Width, g.Height, len(g.Filters)
	a := blas.NewMat(out.X*out.Y, m*n)
	{
		var r int
		for u := 0; u < h.Width; u++ {
			for v := 0; v < h.Height; v++ {
				var s int
				for i := 0; i < g.Width; i++ {
					for j := 0; j < g.Height; j++ {
						a.Set(r, s, f.At(i+u, j+v))
						s++
					}
				}
				r++
			}
		}
	}
	x := blas.NewMat(m*n, k)
	{
		var r int
		for i := 0; i < g.Width; i++ {
			for j := 0; j < g.Height; j++ {
				for p := 0; p < h.Channels; p++ {
					x.Set(r, p, g.Filters[p].At(i, j))
				}
				r++
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

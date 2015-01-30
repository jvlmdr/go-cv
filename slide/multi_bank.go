package slide

import (
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-fftw/fftw"
	"github.com/jvlmdr/lin-go/blas"
)

// MultiBank describes a collection of multi-channel filters.
// All filters must have the same spatial dimension.
type MultiBank struct {
	Width    int
	Height   int
	Channels int
	Filters  []*rimg64.Multi
}

// Size gives the spatial dimension of all filters in the bank.
func (bank *MultiBank) Size() image.Point {
	return image.Pt(bank.Width, bank.Height)
}

// CorrMultiBankNaive computes the correlation of
// a multi-channel image with a bank of multi-channel filters.
// 	h_p[u, v] = sum_q (f_q corr g_pq)[u, v]
func CorrMultiBankNaive(f *rimg64.Multi, g *MultiBank) (*rimg64.Multi, error) {
	out := ValidSize(f.Size(), g.Size())
	if out.X <= 0 || out.Y <= 0 {
		return nil, nil
	}
	h := rimg64.NewMulti(out.X, out.Y, len(g.Filters))
	for u := 0; u < h.Width; u++ {
		for v := 0; v < h.Height; v++ {
			for p := 0; p < h.Channels; p++ {
				var sum float64
				for i := 0; i < g.Width; i++ {
					for j := 0; j < g.Height; j++ {
						for q := 0; q < g.Channels; q++ {
							sum += f.At(i+u, j+v, q) * g.Filters[p].At(i, j, q)
						}
					}
				}
				h.Set(u, v, p, sum)
			}
		}
	}
	return h, nil
}

// CorrMultiBankFFT computes the correlation of
// a multi-channel image with a bank of multi-channel filters.
// 	h_p[u, v] = sum_q (f_q corr g_pq)[u, v]
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
		copyChannelTo(fhat[i], f, i)
		fftw.FFT2To(fhat[i], fhat[i])
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
			// h_p[x] = (G_qp corr F_p)[x]
			// H_p[x] = conj(G_qp[x]) F_p[x]
			addScaleMul(sum, alpha, curr, fhat[q])
		}
		bwd.Execute()
		copyRealToChannel(h, p, sum)
	}
	return h, nil
}

// CorrMultiBankBLAS computes the correlation of
// a multi-channel image with a bank of multi-channel filters.
// 	h_p[u, v] = sum_q (f_q corr g_pq)[u, v]
func CorrMultiBankBLAS(f *rimg64.Multi, g *MultiBank) (*rimg64.Multi, error) {
	out := ValidSize(f.Size(), g.Size())
	if out.X <= 0 || out.Y <= 0 {
		return nil, nil
	}
	// Express as dense matrix multiplication.
	//   h_p[u, v] = sum_q (f_q corr g_pq)[u, v]
	//   Y(h) = A(f) X(g)
	// If the number of input and output channels are Q and P, then
	//   A is (M-m+1)(N-n+1) x mnQ and
	//   X is mnQ x P, so that
	//   Y is (M-m+1)(N-n+1) x P.
	// Note that the time to build the system is therefore
	// affected more by the number of input channels Q than outputs P.

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

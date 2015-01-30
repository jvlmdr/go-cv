package slide

import (
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-fftw/fftw"
	"github.com/jvlmdr/lin-go/blas"
)

// CorrBankStrideNaive computes the strided correlation of
// an image with a bank of filters.
// 	h_p[u, v] = (f corr g_p)[stride*u, stride*v]
func CorrBankStrideNaive(f *rimg64.Image, g *Bank, stride int) (*rimg64.Multi, error) {
	out := ValidSizeStride(f.Size(), g.Size(), stride)
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
						sum += f.At(stride*u+i, stride*v+j) * g.Filters[p].At(i, j)
					}
				}
				h.Set(u, v, p, sum)
			}
		}
	}
	return h, nil
}

// CorrBankStrideFFT computes the strided correlation of
// an image with a bank of filters.
// 	h_p[u, v] = (f corr g_p)[stride*u, stride*v]
func CorrBankStrideFFT(f *rimg64.Image, g *Bank, stride int) (*rimg64.Multi, error) {
	out := ValidSizeStride(f.Size(), g.Size(), stride)
	if out.X <= 0 || out.Y <= 0 {
		return nil, nil
	}
	// Compute strided convolution as the sum over
	// a stride x stride grid of small convolutions.
	grid := image.Pt(stride, stride)
	// But do not divide into a larger grid than the size of the filter.
	// If the filter is smaller than the stride,
	// then some pixels in the image will not affect the output.
	grid.X = min(grid.X, g.Width)
	grid.Y = min(grid.Y, g.Height)
	// Determine the size of the sub-sampled filter.
	gsub := image.Pt(ceilDiv(g.Width, grid.X), ceilDiv(g.Height, grid.Y))
	// The sub-sampled size of the image should be such that
	// the output size is attained.
	fsub := image.Pt(out.X+gsub.X-1, out.Y+gsub.Y-1)

	// Determine optimal size for FFT.
	work, _ := FFT2Size(fsub)
	// Cache FFT of image for convolving with multiple filters.
	// Re-use plan for multiple convolutions too.
	fhat := fftw.NewArray2(work.X, work.Y)
	ffwd := fftw.NewPlan2(fhat, fhat, fftw.Forward, fftw.Estimate)
	defer ffwd.Destroy()
	// FFT for current filter.
	ghat := fftw.NewArray2(work.X, work.Y)
	gfwd := fftw.NewPlan2(ghat, ghat, fftw.Forward, fftw.Estimate)
	defer gfwd.Destroy()
	// Allocate one array per output channel.
	hhat := make([]*fftw.Array2, len(g.Filters))
	for k := range hhat {
		hhat[k] = fftw.NewArray2(work.X, work.Y)
	}
	// Normalization factor.
	alpha := complex(1/float64(work.X*work.Y), 0)
	// Add the convolutions over channels and strides.
	for i := 0; i < grid.X; i++ {
		for j := 0; j < grid.Y; j++ {
			// Take transform of downsampled image given offset (i, j).
			copyStrideTo(fhat, f, stride, image.Pt(i, j))
			ffwd.Execute()
			// Take transform of each downsampled channel given offset (i, j).
			for q := range hhat {
				copyStrideTo(ghat, g.Filters[q], stride, image.Pt(i, j))
				gfwd.Execute()
				addMul(hhat[q], ghat, fhat)
			}
		}
	}
	// Take the inverse transform of each channel.
	h := rimg64.NewMulti(out.X, out.Y, len(g.Filters))
	for q := range hhat {
		scale(alpha, hhat[q])
		fftw.IFFT2To(hhat[q], hhat[q])
		copyRealToChannel(h, q, hhat[q])
	}
	return h, nil
}

// CorrBankStrideBLAS computes the strided correlation of
// an image with a bank of filters.
// 	h_p[u, v] = (f corr g_p)[stride*u, stride*v]
func CorrBankStrideBLAS(f *rimg64.Image, g *Bank, stride int) (*rimg64.Multi, error) {
	out := ValidSizeStride(f.Size(), g.Size(), stride)
	if out.X <= 0 || out.Y <= 0 {
		return nil, nil
	}
	h := rimg64.NewMulti(out.X, out.Y, len(g.Filters))
	// Size of filters.
	m, n := g.Width, g.Height
	// Express as dense matrix multiplication.
	//   h_p[u, v] = sum_q (f_q corr g_pq)[u, v]
	//   h = A(f) X(g)
	// where A is whk by mnk
	// with w = ceil[(M-m+1)/stride],
	//      h = ceil[(N-n+1)/stride].
	a := blas.NewMat(h.Width*h.Height, m*n)
	{
		var r int
		for u := 0; u < h.Width; u++ {
			for v := 0; v < h.Height; v++ {
				var s int
				for i := 0; i < g.Width; i++ {
					for j := 0; j < g.Height; j++ {
						a.Set(r, s, f.At(stride*u+i, stride*v+j))
						s++
					}
				}
				r++
			}
		}
	}
	x := blas.NewMat(m*n, h.Channels)
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

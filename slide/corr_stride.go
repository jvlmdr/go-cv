package slide

import (
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-fftw/fftw"
)

// CorrStrideNaive computes the correlation of
// a multi-channel image with a multi-channel filter.
// 	h[u, v] = sum_q (f_q corr g_q)[u, v]
func CorrStrideNaive(f, g *rimg64.Image, stride int) (*rimg64.Image, error) {
	out := ValidSizeStride(f.Size(), g.Size(), stride)
	h := rimg64.New(out.X, out.Y)
	for i := 0; i < h.Width; i++ {
		for j := 0; j < h.Height; j++ {
			var total float64
			for u := 0; u < g.Width; u++ {
				for v := 0; v < g.Height; v++ {
					p := image.Pt(i, j).Mul(stride).Add(image.Pt(u, v))
					total += f.At(p.X, p.Y) * g.At(u, v)
				}
			}
			h.Set(i, j, total)
		}
	}
	return h, nil
}

// CorrStrideFFT computes the correlation of
// a multi-channel image with a multi-channel filter.
// 	h[u, v] = sum_q (f_q corr g_q)[u, v]
func CorrStrideFFT(f, g *rimg64.Image, stride int) (*rimg64.Image, error) {
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
	// Cache FFT of each channel of image for convolving with multiple filters.
	// Re-use plan for multiple convolutions too.
	fhat := fftw.NewArray2(work.X, work.Y)
	ffwd := fftw.NewPlan2(fhat, fhat, fftw.Forward, fftw.Estimate)
	defer ffwd.Destroy()
	// FFT for current filter.
	curr := fftw.NewArray2(work.X, work.Y)
	gfwd := fftw.NewPlan2(curr, curr, fftw.Forward, fftw.Estimate)
	defer gfwd.Destroy()
	// Normalization factor.
	alpha := complex(1/float64(work.X*work.Y), 0)
	// Add the convolutions over strides.
	hhat := fftw.NewArray2(work.X, work.Y)
	for i := 0; i < grid.X; i++ {
		for j := 0; j < grid.Y; j++ {
			// Copy each downsampled channel and take its transform.
			copyStrideTo(fhat, f, stride, image.Pt(i, j))
			ffwd.Execute()
			copyStrideTo(curr, g, stride, image.Pt(i, j))
			gfwd.Execute()
			addMul(hhat, curr, fhat)
		}
	}
	// Take the inverse transform.
	h := rimg64.New(out.X, out.Y)
	scale(alpha, hhat)
	realIFFTTo(h, hhat)
	return h, nil
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

func realIFFTTo(dst *rimg64.Image, src *fftw.Array2) {
	plan := fftw.NewPlan2(src, src, fftw.Backward, fftw.Estimate)
	defer plan.Destroy()
	plan.Execute()
	for i := 0; i < dst.Width; i++ {
		for j := 0; j < dst.Height; j++ {
			dst.Set(i, j, real(src.At(i, j)))
		}
	}
}

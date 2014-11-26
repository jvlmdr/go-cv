package slide

import (
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-fftw/fftw"
)

// CorrMultiStrideNaive computes the correlation of
// a multi-channel image with a multi-channel filter.
// 	h[u, v] = sum_q (f_q corr g_q)[u, v]
func CorrMultiStrideNaive(f, g *rimg64.Multi, stride int) (*rimg64.Image, error) {
	if err := errIfChannelsNotEq(f, g); err != nil {
		panic(err)
	}
	out := ValidSizeStride(f.Size(), g.Size(), stride)
	h := rimg64.New(out.X, out.Y)
	for i := 0; i < h.Width; i++ {
		for j := 0; j < h.Height; j++ {
			var total float64
			for u := 0; u < g.Width; u++ {
				for v := 0; v < g.Height; v++ {
					p := image.Pt(i, j).Mul(stride).Add(image.Pt(u, v))
					for k := 0; k < f.Channels; k++ {
						total += f.At(p.X, p.Y, k) * g.At(u, v, k)
					}
				}
			}
			h.Set(i, j, total)
		}
	}
	return h, nil
}

// CorrMultiStrideFFT computes the correlation of
// a multi-channel image with a multi-channel filter.
// 	h[u, v] = sum_q (f_q corr g_q)[u, v]
func CorrMultiStrideFFT(f, g *rimg64.Multi, stride int) (*rimg64.Image, error) {
	if err := errIfChannelsNotEq(f, g); err != nil {
		panic(err)
	}
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
	ghat := fftw.NewArray2(work.X, work.Y)
	gfwd := fftw.NewPlan2(ghat, ghat, fftw.Forward, fftw.Estimate)
	defer gfwd.Destroy()
	// Normalization factor.
	alpha := complex(1/float64(work.X*work.Y), 0)
	// Add the convolutions over channels and strides.
	hhat := fftw.NewArray2(work.X, work.Y)
	for k := 0; k < f.Channels; k++ {
		for i := 0; i < grid.X; i++ {
			for j := 0; j < grid.Y; j++ {
				// Copy each downsampled channel and take its transform.
				copyChannelStrideTo(fhat, f, k, stride, image.Pt(i, j))
				ffwd.Execute()
				copyChannelStrideTo(ghat, g, k, stride, image.Pt(i, j))
				gfwd.Execute()
				addMul(hhat, fhat, ghat)
			}
		}
	}
	// Take the inverse transform.
	h := rimg64.New(out.X, out.Y)
	scale(alpha, hhat)
	fftw.IFFT2To(hhat, hhat)
	copyRealTo(h, hhat)
	return h, nil
}

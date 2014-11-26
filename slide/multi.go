package slide

import (
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-fftw/fftw"
	"github.com/jvlmdr/lin-go/blas"
)

// CorrMulti computes the correlation of
// a multi-channel image with a multi-channel filter.
func CorrMulti(f, g *rimg64.Multi) (*rimg64.Image, error) {
	return CorrMultiAlgo(f, g, DefaultAlgo)
}

// CorrMultiAuto computes the correlation of
// a multi-channel image with a multi-channel filter.
// 	h[u, v] = sum_p (f_p corr g_p)[u, v]
// Automatically selects between naive and Fourier-domain convolution.
func CorrMultiAuto(f, g *rimg64.Multi) (*rimg64.Image, error) {
	if err := errIfChannelsNotEq(f, g); err != nil {
		panic(err)
	}
	// Size of output.
	size := ValidSize(f.Size(), g.Size())
	// Return empty image if that's the result.
	if size.Eq(image.ZP) {
		return nil, nil
	}
	// Need to compute one inner product per output element.
	naiveMuls := size.X * size.Y * g.Width * g.Height
	// Optimal FFT size and number of multiplications.
	_, fftMuls := FFT2Size(f.Size())
	// Need to perform two forward and one inverse transform.
	fftMuls *= 3
	// Switch implementation based on image size.
	if fftMuls < naiveMuls {
		return CorrMultiFFT(f, g)
	}
	return CorrMultiNaive(f, g)
}

// CorrMultiNaive computes the correlation of
// a multi-channel image with a multi-channel filter.
// 	h[u, v] = sum_p (f_p corr g_p)[u, v]
func CorrMultiNaive(f, g *rimg64.Multi) (*rimg64.Image, error) {
	if err := errIfChannelsNotEq(f, g); err != nil {
		panic(err)
	}
	out := ValidSize(f.Size(), g.Size())
	if out.Eq(image.ZP) {
		return nil, nil
	}
	h := rimg64.New(out.X, out.Y)
	for i := 0; i < out.X; i++ {
		for j := 0; j < out.Y; j++ {
			var total float64
			for u := 0; u < g.Width; u++ {
				for v := 0; v < g.Height; v++ {
					for p := 0; p < f.Channels; p++ {
						total += f.At(i+u, j+v, p) * g.At(u, v, p)
					}
				}
			}
			h.Set(i, j, total)
		}
	}
	return h, nil
}

// CorrMultiBankFFT computes the correlation of
// a multi-channel image with a multi-channel filter.
// 	h[u, v] = sum_p (f_p corr g_p)[u, v]
func CorrMultiFFT(f, g *rimg64.Multi) (*rimg64.Image, error) {
	if err := errIfChannelsNotEq(f, g); err != nil {
		panic(err)
	}
	out := ValidSize(f.Size(), g.Size())
	if out.Eq(image.ZP) {
		return nil, nil
	}
	work, _ := FFT2Size(f.Size())
	fhat := fftw.NewArray2(work.X, work.Y)
	ghat := fftw.NewArray2(work.X, work.Y)
	ffwd := fftw.NewPlan2(fhat, fhat, fftw.Forward, fftw.Estimate)
	defer ffwd.Destroy()
	gfwd := fftw.NewPlan2(ghat, ghat, fftw.Forward, fftw.Estimate)
	defer gfwd.Destroy()
	hhat := fftw.NewArray2(work.X, work.Y)
	for p := 0; p < f.Channels; p++ {
		// Take transform of each channel.
		copyChannelTo(fhat, f, p)
		ffwd.Execute()
		copyChannelTo(ghat, g, p)
		gfwd.Execute()
		addMul(hhat, ghat, fhat)
	}
	n := float64(work.X * work.Y)
	scale(complex(1/n, 0), hhat)
	fftw.IFFT2To(hhat, hhat)
	h := rimg64.New(out.X, out.Y)
	copyRealTo(h, hhat)
	return h, nil
}

// CorrMultiBLAS computes the correlation of
// a multi-channel image with a multi-channel filter.
// 	h[u, v] = sum_q (f_q corr g_q)[u, v]
func CorrMultiBLAS(f, g *rimg64.Multi) (*rimg64.Image, error) {
	out := ValidSize(f.Size(), g.Size())
	if out.X <= 0 || out.Y <= 0 {
		return nil, nil
	}
	h := rimg64.New(out.X, out.Y)
	// Size of filters.
	m, n, k := g.Width, g.Height, g.Channels
	// Express as dense matrix multiplication.
	//   h[u, v] = sum_q (f_q corr g_q)[u, v]
	//   y(h) = A(f) x(g)
	// where A is (M-m+1)(N-n+1) by mnk.
	a := blas.NewMat(h.Width*h.Height, m*n*k)
	{
		var r int
		for u := 0; u < h.Width; u++ {
			for v := 0; v < h.Height; v++ {
				var s int
				for i := 0; i < g.Width; i++ {
					for j := 0; j < g.Height; j++ {
						for q := 0; q < g.Channels; q++ {
							a.Set(r, s, f.At(u+i, v+j, q))
							s++
						}
					}
				}
				r++
			}
		}
	}
	x := blas.NewMat(m*n*k, 1)
	{
		var r int
		for i := 0; i < g.Width; i++ {
			for j := 0; j < g.Height; j++ {
				for q := 0; q < g.Channels; q++ {
					x.Set(r, 0, g.At(i, j, q))
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
				h.Set(u, v, y.At(r, 0))
				r++
			}
		}
	}
	return h, nil
}

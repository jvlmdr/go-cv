package slide

import (
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-fftw/fftw"
	"github.com/jvlmdr/lin-go/blas"
)

var DefaultAlgo Algo = Auto

// Corr computes the correlation of an image with a filter.
// 	h[u, v] = (f corr g)[u, v]
func Corr(f, g *rimg64.Image) (*rimg64.Image, error) {
	return CorrAlgo(f, g, DefaultAlgo)
}

// CorrAuto computes the correlation of an image with a filter.
// 	h[u, v] = (f corr g)[u, v]
// Automatically selects between naive and Fourier-domain convolution.
func CorrAuto(f, g *rimg64.Image) (*rimg64.Image, error) {
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
		return CorrFFT(f, g)
	}
	return CorrNaive(f, g)
}

// CorrNaive computes the correlation of an image with a filter.
// 	h[u, v] = (f corr g)[u, v]
func CorrNaive(f, g *rimg64.Image) (*rimg64.Image, error) {
	out := ValidSize(f.Size(), g.Size())
	if out.X <= 0 || out.Y <= 0 {
		return nil, nil
	}
	h := rimg64.New(out.X, out.Y)
	for i := 0; i < out.X; i++ {
		for j := 0; j < out.Y; j++ {
			var total float64
			for u := 0; u < g.Width; u++ {
				for v := 0; v < g.Height; v++ {
					total += f.At(i+u, j+v) * g.At(u, v)
				}
			}
			h.Set(i, j, total)
		}
	}
	return h, nil
}

// CorrFFT computes the correlation of an image with a filter.
// 	h[u, v] = (f corr g)[u, v]
func CorrFFT(f, g *rimg64.Image) (*rimg64.Image, error) {
	out := ValidSize(f.Size(), g.Size())
	if out.X <= 0 || out.Y <= 0 {
		return nil, nil
	}
	// Determine optimal size for FFT.
	work, _ := FFT2Size(f.Size())
	fhat := fftw.NewArray2(work.X, work.Y)
	ghat := fftw.NewArray2(work.X, work.Y)
	// Take forward transforms.
	copyImageTo(fhat, f)
	fftw.FFT2To(fhat, fhat)
	copyImageTo(ghat, g)
	fftw.FFT2To(ghat, ghat)
	// Scale such that convolution theorem holds.
	n := float64(work.X * work.Y)
	scaleMul(fhat, complex(1/n, 0), ghat, fhat)
	// Take inverse transform.
	h := rimg64.New(out.X, out.Y)
	fftw.IFFT2To(fhat, fhat)
	copyRealTo(h, fhat)
	return h, nil
}

// CorrBLAS computes the correlation of an image with a filter.
// 	h[u, v] = (f corr g)[u, v]
func CorrBLAS(f, g *rimg64.Image) (*rimg64.Image, error) {
	out := ValidSize(f.Size(), g.Size())
	if out.X <= 0 || out.Y <= 0 {
		return nil, nil
	}
	h := rimg64.New(out.X, out.Y)
	// Size of filters.
	m, n := g.Width, g.Height
	// Express as dense matrix multiplication.
	//   h[u, v] = (f corr g)[u, v]
	//   y(h) = A(f) x(g)
	// where A is (M-m+1)(N-n+1) by mn.
	a := blas.NewMat(h.Width*h.Height, m*n)
	{
		var r int
		for u := 0; u < h.Width; u++ {
			for v := 0; v < h.Height; v++ {
				var s int
				for i := 0; i < g.Width; i++ {
					for j := 0; j < g.Height; j++ {
						a.Set(r, s, f.At(u+i, v+j))
						s++
					}
				}
				r++
			}
		}
	}
	x := blas.NewMat(m*n, 1)
	{
		var r int
		for i := 0; i < g.Width; i++ {
			for j := 0; j < g.Height; j++ {
				x.Set(r, 0, g.At(i, j))
				r++
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

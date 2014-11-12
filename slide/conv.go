package slide

import (
	"image"
	"math/cmplx"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-fftw/fftw"
)

// Conv computes convolution of template g with image f.
// Returns the inner product at all positions such that g lies entirely within f.
//	If h = corr(f, g), then h(t) = sum_{tau} f(t-tau) g(tau).
// Beware: This is typically denoted g * f with the arguments in the opposite order.
//
// Automatically selects between naive and Fourier-domain convolution.
func Conv(f, g *rimg64.Image) *rimg64.Image {
	return convAuto(f, g, false)
}

func Flip(f *rimg64.Image) *rimg64.Image {
	g := rimg64.New(f.Width, f.Height)
	for i := 0; i < f.Width; i++ {
		for j := 0; j < f.Height; j++ {
			g.Set(f.Width-1-i, f.Height-1-j, f.At(i, j))
		}
	}
	return g
}

// Corr computes correlation of template g with image f.
// Returns the inner product at all positions such that g lies entirely within f.
//	If h = corr(f, g), then h(t) = sum_{tau} f(t+tau) g(tau).
// Beware: This is typically denoted g * f with the arguments in the opposite order.
//
// Automatically selects between naive and Fourier-domain convolution.
func Corr(f, g *rimg64.Image) *rimg64.Image {
	return convAuto(f, g, true)
}

func convAuto(f, g *rimg64.Image, corr bool) *rimg64.Image {
	// Size of output.
	size := ValidSize(f.Size(), g.Size())
	// Return empty image if that's the result.
	if size.Eq(image.ZP) {
		return nil
	}
	// Need to compute one inner product per output element.
	naiveMuls := size.X * size.Y * g.Width * g.Height
	// Optimal FFT size and number of multiplications.
	fftSize, fftMuls := FFT2Size(f.Size())
	// Need to perform two forward and one inverse transform.
	fftMuls *= 3
	// Switch implementation based on image size.
	if fftMuls < naiveMuls {
		return convFFT(f, g, fftSize, corr)
	}
	return convNaive(f, g, corr)
}

func convNaive(f, g *rimg64.Image, corr bool) *rimg64.Image {
	r := validRect(f.Size(), g.Size(), corr)
	h := rimg64.New(r.Dx(), r.Dy())
	for i := r.Min.X; i < r.Max.X; i++ {
		for j := r.Min.Y; j < r.Max.Y; j++ {
			var total float64
			for u := 0; u < g.Width; u++ {
				for v := 0; v < g.Height; v++ {
					if corr {
						total += f.At(i+u, j+v) * g.At(u, v)
					} else {
						total += f.At(i-u, j-v) * g.At(u, v)
					}
				}
			}
			h.Set(i-r.Min.X, j-r.Min.Y, total)
		}
	}
	return h
}

// The work parameter specifies the dimension of the FFT.
// The out parameter gives the size of the result.
func convFFT(f, g *rimg64.Image, work image.Point, corr bool) *rimg64.Image {
	x := fftw.NewArray2(work.X, work.Y)
	y := fftw.NewArray2(work.X, work.Y)
	// Copy into FFT arrays.
	copyImageTo(x, f)
	copyImageTo(y, g)
	// Take forward transforms.
	x = fftw.FFT2(x)
	y = fftw.FFT2(y)
	// Multiply in Fourier domain.
	for u := 0; u < work.X; u++ {
		for v := 0; v < work.Y; v++ {
			if corr {
				x.Set(u, v, x.At(u, v)*cmplx.Conj(y.At(u, v)))
			} else {
				x.Set(u, v, x.At(u, v)*y.At(u, v))
			}
		}
	}
	// Take inverse transform.
	x = fftw.IFFT2(x)

	r := validRect(f.Size(), g.Size(), corr)
	// Extract desired region.
	h := rimg64.New(r.Dx(), r.Dy())
	// Scale such that convolution theorem holds.
	n := float64(work.X * work.Y)
	for u := r.Min.X; u < r.Max.X; u++ {
		for v := r.Min.Y; v < r.Max.Y; v++ {
			h.Set(u-r.Min.X, v-r.Min.Y, real(x.At(u, v))/n)
		}
	}
	return h
}

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

func ConvNaive(f, g *rimg64.Image) *rimg64.Image {
	out := ValidSize(f.Size(), g.Size())
	if out.Eq(image.ZP) {
		return nil
	}
	return convNaive(f, g, false)
}

func CorrNaive(f, g *rimg64.Image) *rimg64.Image {
	out := ValidSize(f.Size(), g.Size())
	if out.Eq(image.ZP) {
		return nil
	}
	return convNaive(f, g, true)
}

func ConvFFT(f, g *rimg64.Image) *rimg64.Image {
	out := ValidSize(f.Size(), g.Size())
	if out.Eq(image.ZP) {
		return nil
	}
	work, _ := FFT2Size(f.Size())
	return convFFT(f, g, work, false)
}

func CorrFFT(f, g *rimg64.Image) *rimg64.Image {
	out := ValidSize(f.Size(), g.Size())
	if out.Eq(image.ZP) {
		return nil
	}
	work, _ := FFT2Size(f.Size())
	return convFFT(f, g, work, true)
}

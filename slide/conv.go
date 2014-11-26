package slide

import (
	"image"
	"math/cmplx"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-fftw/fftw"
)

// Flip mirrors an image in both dimensions.
func Flip(f *rimg64.Image) *rimg64.Image {
	g := rimg64.New(f.Width, f.Height)
	for i := 0; i < f.Width; i++ {
		for j := 0; j < f.Height; j++ {
			g.Set(f.Width-1-i, f.Height-1-j, f.At(i, j))
		}
	}
	return g
}

// Conv computes convolution of template g with image f.
// Returns the inner product at all positions such that g lies entirely within f.
//	If h = corr(f, g), then h(t) = sum_{tau} f(t-tau) g(tau).
// Beware: This is typically denoted g * f with the arguments in the opposite order.
//
// Automatically selects between naive and Fourier-domain convolution.
func Conv(f, g *rimg64.Image) (*rimg64.Image, error) {
	return convAuto(f, g, false)
}

// Corr computes correlation of template g with image f.
// Returns the inner product at all positions such that g lies entirely within f.
//	If h = corr(f, g), then h(t) = sum_{tau} f(t+tau) g(tau).
// Beware: This is typically denoted g * f with the arguments in the opposite order.
//
// Automatically selects between naive and Fourier-domain convolution.
func Corr(f, g *rimg64.Image) (*rimg64.Image, error) {
	return convAuto(f, g, true)
}

func convAuto(f, g *rimg64.Image, corr bool) (*rimg64.Image, error) {
	// Size of output.
	size := ValidSize(f.Size(), g.Size())
	// Return empty image if that's the result.
	if size.Eq(image.ZP) {
		return nil, nil
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

func convNaive(f, g *rimg64.Image, corr bool) (*rimg64.Image, error) {
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
	return h, nil
}

// The work parameter specifies the dimension of the FFT.
// The out parameter gives the size of the result.
func convFFT(f, g *rimg64.Image, work image.Point, corr bool) (*rimg64.Image, error) {
	fhat := fftw.NewArray2(work.X, work.Y)
	ghat := fftw.NewArray2(work.X, work.Y)
	// Copy into FFT arrays.
	copyImageTo(fhat, f)
	copyImageTo(ghat, g)
	// Take forward transforms.
	fftw.FFT2To(fhat, fhat)
	fftw.FFT2To(ghat, ghat)
	// Multiply in Fourier domain.
	for u := 0; u < work.X; u++ {
		for v := 0; v < work.Y; v++ {
			if corr {
				fhat.Set(u, v, fhat.At(u, v)*cmplx.Conj(ghat.At(u, v)))
			} else {
				fhat.Set(u, v, fhat.At(u, v)*ghat.At(u, v))
			}
		}
	}
	// Take inverse transform.
	fftw.IFFT2To(fhat, fhat)

	r := validRect(f.Size(), g.Size(), corr)
	// Extract desired region.
	h := rimg64.New(r.Dx(), r.Dy())
	// Scale such that convolution theorem holds.
	n := float64(work.X * work.Y)
	for u := r.Min.X; u < r.Max.X; u++ {
		for v := r.Min.Y; v < r.Max.Y; v++ {
			h.Set(u-r.Min.X, v-r.Min.Y, real(fhat.At(u, v))/n)
		}
	}
	return h, nil
}

func ConvNaive(f, g *rimg64.Image) (*rimg64.Image, error) {
	out := ValidSize(f.Size(), g.Size())
	if out.Eq(image.ZP) {
		return nil, nil
	}
	return convNaive(f, g, false)
}

func CorrNaive(f, g *rimg64.Image) (*rimg64.Image, error) {
	out := ValidSize(f.Size(), g.Size())
	if out.Eq(image.ZP) {
		return nil, nil
	}
	return convNaive(f, g, true)
}

func ConvFFT(f, g *rimg64.Image) (*rimg64.Image, error) {
	out := ValidSize(f.Size(), g.Size())
	if out.Eq(image.ZP) {
		return nil, nil
	}
	work, _ := FFT2Size(f.Size())
	return convFFT(f, g, work, false)
}

func CorrFFT(f, g *rimg64.Image) (*rimg64.Image, error) {
	out := ValidSize(f.Size(), g.Size())
	if out.Eq(image.ZP) {
		return nil, nil
	}
	work, _ := FFT2Size(f.Size())
	return convFFT(f, g, work, true)
}

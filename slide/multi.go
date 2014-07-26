package slide

import (
	"fmt"
	"image"
	"math/cmplx"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-fftw/fftw"
)

// Performs correlation of multi-channel images.
// Returns sum over channels.
func CorrMulti(f, g *rimg64.Multi) *rimg64.Image {
	if f.Channels != g.Channels {
		err := fmt.Errorf("different number of channels: %d, %d", f.Channels, g.Channels)
		panic(err)
	}

	size := outputSize(f.Size(), g.Size())
	// Return empty image if that's the result.
	if size.X == 0 || size.Y == 0 {
		return nil
	}
	// Switch implementation based on image size.
	if useFourier(f.Size(), g.Size()) {
		return corrMultiFFT(f, g)
	}
	return corrMultiNaive(f, g)
}

func corrMultiNaive(f, g *rimg64.Multi) *rimg64.Image {
	size := outputSize(f.Size(), g.Size())
	// Return empty image if that's the result.
	if size.Eq(image.ZP) {
		return nil
	}

	h := rimg64.New(size.X, size.Y)
	for i := 0; i < size.X; i++ {
		for j := 0; j < size.Y; j++ {
			var total float64
			for k := 0; k < f.Channels; k++ {
				for u := 0; u < g.Width; u++ {
					for v := 0; v < g.Height; v++ {
						total += f.At(i+u, j+v, k) * g.At(u, v, k)
					}
				}
			}
			h.Set(i, j, total)
		}
	}
	return h
}

func corrMultiFFT(f, g *rimg64.Multi) *rimg64.Image {
	size := outputSize(f.Size(), g.Size())
	// Return empty image if that's the result.
	if size.Eq(image.ZP) {
		return nil
	}

	// Scale such that convolution theorem holds.
	n := float64(f.Width) * float64(f.Height)

	fhat := fftw.NewArray2(f.Width, f.Height)
	ghat := fftw.NewArray2(f.Width, f.Height)

	h := rimg64.New(size.X, size.Y)
	for k := 0; k < f.Channels; k++ {
		// Copy into FFT arrays.
		for u := 0; u < f.Width; u++ {
			for v := 0; v < f.Height; v++ {
				fhat.Set(u, v, complex(f.At(u, v, k), 0))
				if u < g.Width && v < g.Height {
					ghat.Set(u, v, complex(g.At(u, v, k), 0))
				} else {
					ghat.Set(u, v, 0)
				}
			}
		}
		// Take forward transforms.
		fhat = fftw.FFT2(fhat)
		ghat = fftw.FFT2(ghat)
		// Multiply in Fourier domain.
		for u := 0; u < f.Width; u++ {
			for v := 0; v < f.Height; v++ {
				fhat.Set(u, v, fhat.At(u, v)*cmplx.Conj(ghat.At(u, v)))
			}
		}
		// Take inverse transform.
		fhat = fftw.IFFT2(fhat)

		// Sum response over multiple channels.
		// Scale such that convolution theorem holds.
		for u := 0; u < size.X; u++ {
			for v := 0; v < size.Y; v++ {
				h.Set(u, v, h.At(u, v)+real(fhat.At(u, v))/n)
			}
		}
	}
	return h
}

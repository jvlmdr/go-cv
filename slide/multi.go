package slide

import (
	"github.com/jackvalmadre/go-cv/rimg64"
	"github.com/jackvalmadre/go-fftw/fftw"

	"fmt"
	"image"
	"math/cmplx"
)

// Takes inner product of g with f at all positions such that it lies entirely within f.
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

	x := fftw.NewArray2(f.Width, f.Height)
	y := fftw.NewArray2(f.Width, f.Height)

	h := rimg64.New(size.X, size.Y)
	for k := 0; k < f.Channels; k++ {
		// Copy into FFT arrays.
		for u := 0; u < f.Width; u++ {
			for v := 0; v < f.Height; v++ {
				x.Set(u, v, complex(f.At(u, v, k), 0))
				if u < g.Width && v < g.Height {
					y.Set(u, v, complex(g.At(u, v, k), 0))
				} else {
					y.Set(u, v, 0)
				}
			}
		}
		// Take forward transforms.
		x = fftw.FFT2(x)
		y = fftw.FFT2(y)
		// Multiply in Fourier domain.
		for u := 0; u < f.Width; u++ {
			for v := 0; v < f.Height; v++ {
				x.Set(u, v, x.At(u, v)*cmplx.Conj(y.At(u, v)))
			}
		}
		// Take inverse transform.
		x = fftw.IFFT2(x)

		// Sum response over multiple channels.
		// Scale such that convolution theorem holds.
		for u := 0; u < size.X; u++ {
			for v := 0; v < size.Y; v++ {
				h.Set(u, v, h.At(u, v)+real(x.At(u, v))/n)
			}
		}
	}
	return h
}

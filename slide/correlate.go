package slide

import (
	"github.com/jackvalmadre/go-cv"
	"github.com/jackvalmadre/go-fftw"

	"image"
	"math/cmplx"
)

// Returns the number of positions such that the template g lies entirely inside the image f.
func outputSize(f, g image.Point) image.Point {
	var h image.Point
	h.X = max(f.X-g.X+1, 0)
	h.Y = max(f.Y-g.Y+1, 0)
	return h
}

// Computes approximate expense of each approach, return true if Fourier is better.
func useFourier(f, g image.Point) bool {
	h := outputSize(f, g)
	// One dot product per output pixel.
	naive := h.X * h.Y * g.X * g.Y
	// Two forward transforms and an inverse transform.
	fourier := 3 * f.X * f.Y * logb(f.X*f.Y)
	return fourier < naive
}

// Computes correlation of template g with image f.
//
// Takes inner product of g with f at all positions such that it lies entirely within f.
//
// Automatically selects between naive and Fourier-domain convolution.
func CorrelateImages(f, g cv.RealImage) cv.RealImage {
	size := outputSize(f.Size(), g.Size())
	// Return empty image if that's the result.
	if size.Eq(image.ZP) {
		return cv.RealImage{}
	}
	// Switch implementation based on image size.
	if useFourier(f.Size(), g.Size()) {
		return correlateImagesFourier(f, g)
	}
	return correlateImagesNaive(f, g)
}

func correlateImagesNaive(f, g cv.RealImage) cv.RealImage {
	size := outputSize(f.Size(), g.Size())
	// Return empty image if that's the result.
	if size.Eq(image.ZP) {
		return cv.RealImage{}
	}

	h := cv.NewRealImage(size.X, size.Y)
	for i := 0; i < size.X; i++ {
		for j := 0; j < size.Y; j++ {
			var total float64
			for u := 0; u < g.Width; u++ {
				for v := 0; v < g.Height; v++ {
					total += f.At(i+u, j+v) * g.At(u, v)
				}
			}
			h.Set(i, j, total)
		}
	}
	return h
}

func correlateImagesFourier(f, g cv.RealImage) cv.RealImage {
	size := outputSize(f.Size(), g.Size())
	// Return empty image if that's the result.
	if size.Eq(image.ZP) {
		return cv.RealImage{}
	}

	x := fftw.Alloc2d(f.Width, f.Height)
	defer fftw.Free2d(x)
	y := fftw.Alloc2d(f.Width, f.Height)
	defer fftw.Free2d(y)
	// Copy into FFT arrays.
	for u := 0; u < f.Width; u++ {
		for v := 0; v < f.Height; v++ {
			x[u][v] = complex(f.At(u, v), 0)
		}
	}
	for u := 0; u < g.Width; u++ {
		for v := 0; v < g.Height; v++ {
			y[u][v] = complex(g.At(u, v), 0)
		}
	}

	// Take forward transforms.
	fftw.Dft2d(x, x, fftw.Forward, fftw.Estimate)
	fftw.Dft2d(y, y, fftw.Forward, fftw.Estimate)
	// Multiply in Fourier domain.
	for u := 0; u < f.Width; u++ {
		for v := 0; v < f.Height; v++ {
			x[u][v] = x[u][v] * cmplx.Conj(y[u][v])
		}
	}
	// Take inverse transform.
	fftw.Dft2d(x, x, fftw.Backward, fftw.Estimate)

	// Extract desired region.
	h := cv.NewRealImage(size.X, size.Y)
	// Scale such that convolution theorem holds.
	n := float64(f.Width) * float64(f.Height)
	for u := 0; u < size.X; u++ {
		for v := 0; v < size.Y; v++ {
			h.Set(u, v, real(x[u][v])/n)
		}
	}
	return h
}

// Takes inner product of g with f at all positions such that it lies entirely within f.
func CorrelateVectorImages(f, g cv.RealVectorImage) cv.RealImage {
	if f.Channels != g.Channels {
		panic("Number of channels does not match")
	}
	size := outputSize(f.ImageSize(), g.ImageSize())
	// Return empty image if that's the result.
	if size.X == 0 || size.Y == 0 {
		return cv.RealImage{nil, size.X, size.Y}
	}
	// Switch implementation based on image size.
	if useFourier(f.ImageSize(), g.ImageSize()) {
		return correlateVectorImagesFourier(f, g)
	}
	return correlateVectorImagesNaive(f, g)
}

func correlateVectorImagesNaive(f, g cv.RealVectorImage) cv.RealImage {
	size := outputSize(f.ImageSize(), g.ImageSize())
	// Return empty image if that's the result.
	if size.Eq(image.ZP) {
		return cv.RealImage{}
	}

	h := cv.NewRealImage(size.X, size.Y)
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

func correlateVectorImagesFourier(f, g cv.RealVectorImage) cv.RealImage {
	size := outputSize(f.ImageSize(), g.ImageSize())
	// Return empty image if that's the result.
	if size.Eq(image.ZP) {
		return cv.RealImage{}
	}

	// Scale such that convolution theorem holds.
	n := float64(f.Width) * float64(f.Height)

	x := fftw.Alloc2d(f.Width, f.Height)
	defer fftw.Free2d(x)
	y := fftw.Alloc2d(f.Width, f.Height)
	defer fftw.Free2d(y)

	h := cv.NewRealImage(size.X, size.Y)
	for k := 0; k < f.Channels; k++ {
		// Copy into FFT arrays.
		for u := 0; u < f.Width; u++ {
			for v := 0; v < f.Height; v++ {
				x[u][v] = complex(f.At(u, v, k), 0)
				if u < g.Width && v < g.Height {
					y[u][v] = complex(g.At(u, v, k), 0)
				} else {
					y[u][v] = 0
				}
			}
		}
		// Take forward transforms.
		fftw.Dft2d(x, x, fftw.Forward, fftw.Estimate)
		fftw.Dft2d(y, y, fftw.Forward, fftw.Estimate)
		// Multiply in Fourier domain.
		for u := 0; u < f.Width; u++ {
			for v := 0; v < f.Height; v++ {
				x[u][v] = x[u][v] * cmplx.Conj(y[u][v])
			}
		}
		// Take inverse transform.
		fftw.Dft2d(x, x, fftw.Backward, fftw.Estimate)

		// Sum response over multiple channels.
		// Scale such that convolution theorem holds.
		for u := 0; u < size.X; u++ {
			for v := 0; v < size.Y; v++ {
				h.Set(u, v, h.At(u, v)+real(x[u][v])/n)
			}
		}
	}
	return h
}

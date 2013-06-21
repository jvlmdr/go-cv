package cv

import (
	"github.com/jackvalmadre/go-fftw"
	"image"
	"math/cmplx"
)

// Returns the number of positions such that the template lies entirely inside the image.
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

// Takes inner product of g with f at all positions such that it lies entirely within f.
func CorrelateImages(f, g RealImage) RealImage {
	size := outputSize(f.Size(), g.Size())
	// Return empty image if that's the result.
	if size.X == 0 || size.Y == 0 {
		return RealImage{nil, size.X, size.Y}
	}
	// Switch implementation based on image size.
	if useFourier(f.Size(), g.Size()) {
		return correlateImagesFourier(f, g)
	}
	return correlateImagesNaive(f, g)
}

func correlateImagesNaive(f, g RealImage) RealImage {
	size := outputSize(f.Size(), g.Size())
	h := NewRealImage(size.X, size.Y)
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

func correlateImagesFourier(f, g RealImage) RealImage {
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
	// Multiply and scale such that convolution theorem holds.
	n := float64(f.Width) * float64(f.Height)
	for u := 0; u < f.Width; u++ {
		for v := 0; v < f.Height; v++ {
			x[u][v] = complex(1/n, 0) * cmplx.Conj(x[u][v]) * y[u][v]
		}
	}
	// Take inverse transform.
	fftw.Dft2d(x, x, fftw.Backward, fftw.Estimate)
	// Extract desired region.
	size := outputSize(f.Size(), g.Size())
	h := NewRealImage(size.X, size.Y)
	for u := 0; u < size.X; u++ {
		for v := 0; v < size.Y; v++ {
			h.Set(u, v, real(x[u][v]))
		}
	}
	return h
}

// Takes inner product of g with f at all positions such that it lies entirely within f.
func CorrelateVectorImages(f, g RealVectorImage) RealImage {
	if f.Channels != g.Channels {
		panic("Number of channels does not match")
	}
	size := outputSize(f.ImageSize(), g.ImageSize())
	// Return empty image if that's the result.
	if size.X == 0 || size.Y == 0 {
		return RealImage{nil, size.X, size.Y}
	}
	// Switch implementation based on image size.
	if useFourier(f.ImageSize(), g.ImageSize()) {
		return correlateVectorImagesFourier(f, g)
	}
	return correlateVectorImagesNaive(f, g)
}

func correlateVectorImagesNaive(f, g RealVectorImage) RealImage {
	size := outputSize(f.ImageSize(), g.ImageSize())
	h := NewRealImage(size.X, size.Y)
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

func correlateVectorImagesFourier(f, g RealVectorImage) RealImage {
	x := fftw.Alloc2d(f.Width, f.Height)
	defer fftw.Free2d(x)
	y := fftw.Alloc2d(f.Width, f.Height)
	defer fftw.Free2d(y)

	size := outputSize(f.ImageSize(), g.ImageSize())
	h := NewRealImage(size.X, size.Y)
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
		// Multiply and scale such that convolution theorem holds.
		n := float64(f.Width) * float64(f.Height)
		for u := 0; u < f.Width; u++ {
			for v := 0; v < f.Height; v++ {
				x[u][v] = complex(1/n, 0) * cmplx.Conj(x[u][v]) * y[u][v]
			}
		}
		// Take inverse transform.
		fftw.Dft2d(x, x, fftw.Backward, fftw.Estimate)
		// Sum response over multiple channels.
		for u := 0; u < size.X; u++ {
			for v := 0; v < size.Y; v++ {
				h.Set(u, v, h.At(u, v)+real(x[u][v]))
			}
		}
	}
	return h
}

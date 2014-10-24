package slide

import (
	"fmt"
	"image"
	"math/cmplx"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-fftw/fftw"
)

func ConvMulti(f, g *rimg64.Multi) *rimg64.Image {
	return convMulti(f, g, false)
}

func CorrMulti(f, g *rimg64.Multi) *rimg64.Image {
	return convMulti(f, g, true)
}

func errIfChannelsNotEq(f, g *rimg64.Multi) error {
	if f.Channels != g.Channels {
		return fmt.Errorf("different number of channels: %d, %d", f.Channels, g.Channels)
	}
	return nil
}

// Performs correlation of multi-channel images.
// Returns sum over channels.
func convMulti(f, g *rimg64.Multi, corr bool) *rimg64.Image {
	if err := errIfChannelsNotEq(f, g); err != nil {
		panic(err)
	}
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
		return convMultiFFT(f, g, fftSize, corr)
	}
	return convMultiNaive(f, g, corr)
}

func convMultiNaive(f, g *rimg64.Multi, corr bool) *rimg64.Image {
	r := validRect(f.Size(), g.Size(), corr)
	h := rimg64.New(r.Dx(), r.Dy())
	for i := r.Min.X; i < r.Max.X; i++ {
		for j := r.Min.Y; j < r.Max.Y; j++ {
			var total float64
			for p := 0; p < f.Channels; p++ {
				for u := 0; u < g.Width; u++ {
					for v := 0; v < g.Height; v++ {
						if corr {
							total += f.At(i+u, j+v, p) * g.At(u, v, p)
						} else {
							total += f.At(i-u, j-v, p) * g.At(u, v, p)
						}
					}
				}
			}
			h.Set(i-r.Min.X, j-r.Min.Y, total)
		}
	}
	return h
}

func convMultiFFT(f, g *rimg64.Multi, work image.Point, corr bool) *rimg64.Image {
	x := fftw.NewArray2(work.X, work.Y)
	y := fftw.NewArray2(work.X, work.Y)
	fftX := fftw.NewPlan2(x, x, fftw.Forward, fftw.Estimate)
	defer fftX.Destroy()
	fftY := fftw.NewPlan2(y, y, fftw.Forward, fftw.Estimate)
	defer fftY.Destroy()
	ifftX := fftw.NewPlan2(x, x, fftw.Backward, fftw.Estimate)
	defer ifftX.Destroy()

	r := validRect(f.Size(), g.Size(), corr)
	h := rimg64.New(r.Dx(), r.Dy())
	n := float64(work.X * work.Y)

	for p := 0; p < f.Channels; p++ {
		copyChannelTo(x, f, p)
		copyChannelTo(y, g, p)
		fftX.Execute()
		fftY.Execute()
		for u := 0; u < work.X; u++ {
			for v := 0; v < work.Y; v++ {
				if corr {
					x.Set(u, v, x.At(u, v)*cmplx.Conj(y.At(u, v)))
				} else {
					x.Set(u, v, x.At(u, v)*y.At(u, v))
				}
			}
		}
		ifftX.Execute()
		// Sum response over multiple channels.
		// Scale such that convolution theorem holds.
		for u := r.Min.X; u < r.Max.X; u++ {
			for v := r.Min.Y; v < r.Max.Y; v++ {
				i, j := u-r.Min.X, v-r.Min.Y
				h.Set(i, j, h.At(i, j)+real(x.At(u, v))/n)
			}
		}
	}
	return h
}

func copyChannelTo(x *fftw.Array2, f *rimg64.Multi, p int) {
	w, h := x.Dims()
	for u := 0; u < w; u++ {
		for v := 0; v < h; v++ {
			if u < f.Width && v < f.Height {
				x.Set(u, v, complex(f.At(u, v, p), 0))
			} else {
				x.Set(u, v, 0)
			}
		}
	}
}

func ConvMultiNaive(f, g *rimg64.Multi) *rimg64.Image {
	if err := errIfChannelsNotEq(f, g); err != nil {
		panic(err)
	}
	out := ValidSize(f.Size(), g.Size())
	if out.Eq(image.ZP) {
		return nil
	}
	return convMultiNaive(f, g, false)
}

func CorrMultiNaive(f, g *rimg64.Multi) *rimg64.Image {
	if err := errIfChannelsNotEq(f, g); err != nil {
		panic(err)
	}
	out := ValidSize(f.Size(), g.Size())
	if out.Eq(image.ZP) {
		return nil
	}
	return convMultiNaive(f, g, true)
}

func ConvMultiFFT(f, g *rimg64.Multi) *rimg64.Image {
	if err := errIfChannelsNotEq(f, g); err != nil {
		panic(err)
	}
	out := ValidSize(f.Size(), g.Size())
	if out.Eq(image.ZP) {
		return nil
	}
	work, _ := FFT2Size(f.Size())
	return convMultiFFT(f, g, work, false)
}

func CorrMultiFFT(f, g *rimg64.Multi) *rimg64.Image {
	if err := errIfChannelsNotEq(f, g); err != nil {
		panic(err)
	}
	out := ValidSize(f.Size(), g.Size())
	if out.Eq(image.ZP) {
		return nil
	}
	work, _ := FFT2Size(f.Size())
	return convMultiFFT(f, g, work, true)
}

package slide

import (
	"image"
	"math/cmplx"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-fftw/fftw"
)

func ConvMulti(f, g *rimg64.Multi) (*rimg64.Image, error) {
	return convMultiAuto(f, g, false)
}

func CorrMulti(f, g *rimg64.Multi) (*rimg64.Image, error) {
	return convMultiAuto(f, g, true)
}

// Performs correlation of multi-channel images.
// Returns sum over channels.
func convMultiAuto(f, g *rimg64.Multi, corr bool) (*rimg64.Image, error) {
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
	fftSize, fftMuls := FFT2Size(f.Size())
	// Need to perform two forward and one inverse transform.
	fftMuls *= 3
	// Switch implementation based on image size.
	if fftMuls < naiveMuls {
		return convMultiFFT(f, g, fftSize, corr)
	}
	return convMultiNaive(f, g, corr)
}

func convMultiNaive(f, g *rimg64.Multi, corr bool) (*rimg64.Image, error) {
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
	return h, nil
}

func convMultiFFT(f, g *rimg64.Multi, work image.Point, corr bool) (*rimg64.Image, error) {
	fhat := fftw.NewArray2(work.X, work.Y)
	ghat := fftw.NewArray2(work.X, work.Y)
	ffwd := fftw.NewPlan2(fhat, fhat, fftw.Forward, fftw.Estimate)
	defer ffwd.Destroy()
	gfwd := fftw.NewPlan2(ghat, ghat, fftw.Forward, fftw.Estimate)
	defer gfwd.Destroy()
	finv := fftw.NewPlan2(fhat, fhat, fftw.Backward, fftw.Estimate)
	defer finv.Destroy()

	r := validRect(f.Size(), g.Size(), corr)
	h := rimg64.New(r.Dx(), r.Dy())
	n := float64(work.X * work.Y)

	for p := 0; p < f.Channels; p++ {
		copyChannelTo(fhat, f, p)
		copyChannelTo(ghat, g, p)
		ffwd.Execute()
		gfwd.Execute()
		for u := 0; u < work.X; u++ {
			for v := 0; v < work.Y; v++ {
				if corr {
					fhat.Set(u, v, fhat.At(u, v)*cmplx.Conj(ghat.At(u, v)))
				} else {
					fhat.Set(u, v, fhat.At(u, v)*ghat.At(u, v))
				}
			}
		}
		finv.Execute()
		// Sum response over multiple channels.
		// Scale such that convolution theorem holds.
		for u := r.Min.X; u < r.Max.X; u++ {
			for v := r.Min.Y; v < r.Max.Y; v++ {
				i, j := u-r.Min.X, v-r.Min.Y
				h.Set(i, j, h.At(i, j)+real(fhat.At(u, v))/n)
			}
		}
	}
	return h, nil
}

func ConvMultiNaive(f, g *rimg64.Multi) (*rimg64.Image, error) {
	if err := errIfChannelsNotEq(f, g); err != nil {
		panic(err)
	}
	out := ValidSize(f.Size(), g.Size())
	if out.Eq(image.ZP) {
		return nil, nil
	}
	return convMultiNaive(f, g, false)
}

func CorrMultiNaive(f, g *rimg64.Multi) (*rimg64.Image, error) {
	if err := errIfChannelsNotEq(f, g); err != nil {
		panic(err)
	}
	out := ValidSize(f.Size(), g.Size())
	if out.Eq(image.ZP) {
		return nil, nil
	}
	return convMultiNaive(f, g, true)
}

func ConvMultiFFT(f, g *rimg64.Multi) (*rimg64.Image, error) {
	if err := errIfChannelsNotEq(f, g); err != nil {
		panic(err)
	}
	out := ValidSize(f.Size(), g.Size())
	if out.Eq(image.ZP) {
		return nil, nil
	}
	work, _ := FFT2Size(f.Size())
	return convMultiFFT(f, g, work, false)
}

func CorrMultiFFT(f, g *rimg64.Multi) (*rimg64.Image, error) {
	if err := errIfChannelsNotEq(f, g); err != nil {
		panic(err)
	}
	out := ValidSize(f.Size(), g.Size())
	if out.Eq(image.ZP) {
		return nil, nil
	}
	work, _ := FFT2Size(f.Size())
	return convMultiFFT(f, g, work, true)
}

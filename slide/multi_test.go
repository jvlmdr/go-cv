package slide_test

import (
	"image"
	"math"
	"testing"

	"github.com/jvlmdr/go-cv/slide"
)

// Compare naive and Fourier implementations.
func TestCorrMulti_FFTVsNaive(t *testing.T) {
	const (
		m   = 40
		n   = 30
		w   = 100
		h   = 80
		c   = 8
		eps = 1e-9
	)

	f := randMulti(w, h, c)
	g := randMulti(m, n, c)

	naive := slide.CorrMultiNaive(f, g)
	fourier := slide.CorrMultiFFT(f, g)

	if !naive.Size().Eq(fourier.Size()) {
		t.Fatalf("size mismatch (naive %v, fourier %v)", naive.Size(), fourier.Size())
	}

	for x := 0; x < naive.Width; x++ {
		for y := 0; y < naive.Height; y++ {
			xy := image.Pt(x, y)
			if math.Abs(naive.At(x, y)-fourier.At(x, y)) > eps {
				t.Errorf("value mismatch at %v (naive %g, fourier %g)", xy, naive.At(x, y), fourier.At(x, y))
			}
		}
	}
}

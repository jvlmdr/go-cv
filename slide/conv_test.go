package slide_test

import (
	"image"
	"math"
	"testing"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-cv/slide"
)

func TestCorr(t *testing.T) {
	const eps = 1e-9
	f := rimg64.FromRows([][]float64{
		{1, 2, 3, 4, 5},
		{2, 5, 4, 1, 3},
		{5, 4, 3, 2, 1},
	})
	g := rimg64.FromRows([][]float64{
		{3, 1, 5},
		{2, 4, 1},
	})
	cases := []struct {
		I, J int
		Want float64
	}{
		// <g, [1 2 3; 2 5 4]>
		{0, 0, 3*1 + 1*2 + 5*3 + 2*2 + 4*5 + 1*4},
		// <g, [2 3 4; 5 4 1]>
		{1, 0, 3*2 + 1*3 + 5*4 + 2*5 + 4*4 + 1*1},
		// <g, [3 4 5; 4 1 3]>
		{2, 0, 3*3 + 1*4 + 5*5 + 2*4 + 4*1 + 1*3},
		// <g, [2 5 4; 5 4 3]>
		{0, 1, 3*2 + 1*5 + 5*4 + 2*5 + 4*4 + 1*3},
		// <g, [5 4 1; 4 3 2]>
		{1, 1, 3*5 + 1*4 + 5*1 + 2*4 + 4*3 + 1*2},
		// <g, [4 1 3; 3 2 1]>
		{2, 1, 3*4 + 1*1 + 5*3 + 2*3 + 4*2 + 1*1},
	}

	h := slide.Corr(f, g)
	if h.Width != 3 || h.Height != 2 {
		t.Fatalf("wrong size: want %dx%d, got %dx%d", 3, 2, h.Width, h.Height)
	}
	for _, c := range cases {
		if got := h.At(c.I, c.J); math.Abs(got-c.Want) > eps {
			t.Errorf(
				"not equal: (i, j) = (%d, %d): want %.5g, got %.5g",
				c.I, c.J, c.Want, got,
			)
		}
	}
}

func TestConv_vsFlipCorr(t *testing.T) {
	const (
		eps  = 1e-9
		M, N = 601, 599
		m, n = 32, 64
	)
	f := randImage(M, N)
	g := randImage(m, n)
	// Flip g to obtain h.
	h := rimg64.New(m, n)
	for u := 0; u < m; u++ {
		for v := 0; v < n; v++ {
			h.Set(u, v, g.At(m-1-u, n-1-v))
		}
	}
	gConvF := slide.Conv(f, g)
	gCorrF := slide.Corr(f, g)
	hConvF := slide.Conv(f, h)
	hCorrF := slide.Corr(f, h)
	checkImageEq(t, gCorrF, hConvF, eps)
	checkImageEq(t, hCorrF, gConvF, eps)
}

func checkImageEq(t *testing.T, want, got *rimg64.Image, eps float64) {
	if !want.Size().Eq(got.Size()) {
		t.Errorf("different size: want %v, got %v", want.Size(), got.Size())
		return
	}
	for i := 0; i < want.Width; i++ {
		for j := 0; j < want.Height; j++ {
			a, b := want.At(i, j), got.At(i, j)
			if math.Abs(a-b) > eps {
				t.Errorf("different at %d, %d: want %g, got %g", i, j, a, b)
			}
		}
	}
}

func BenchmarkCorrFFT(b *testing.B) {
	const M, N = 800, 600
	const m, n = 20, 40
	for i := 0; i < b.N; i++ {
		x := randImage(M, N)
		y := randImage(m, n)
		slide.CorrFFT(x, y)
	}
}

func BenchmarkCorrNaive(b *testing.B) {
	const M, N = 800, 600
	const m, n = 20, 40
	for i := 0; i < b.N; i++ {
		x := randImage(M, N)
		y := randImage(m, n)
		slide.CorrNaive(x, y)
	}
}

// Compare naive and Fourier implementations.
func TestCorr_fftVsNaive(t *testing.T) {
	const (
		m   = 40
		n   = 30
		w   = 100
		h   = 80
		eps = 1e-12
	)

	f := randImage(w, h)
	g := randImage(m, n)

	naive := slide.CorrNaive(f, g)
	fourier := slide.CorrFFT(f, g)

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

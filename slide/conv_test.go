package slide_test

import (
	"math"
	"testing"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-cv/slide"
)

func TestCorrNaive(t *testing.T) {
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

	h, err := slide.CorrNaive(f, g)
	if err != nil {
		t.Fatal(err)
	}
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
		eps  = 1e-6
		M, N = 601, 599
		m, n = 32, 64
	)
	f := randImage(M, N)
	g := randImage(m, n)
	// Flip g to obtain h.
	h := slide.Flip(g)
	gConvF, err := slide.Conv(f, g)
	if err != nil {
		t.Fatal(err)
	}
	gCorrF, err := slide.Corr(f, g)
	if err != nil {
		t.Fatal(err)
	}
	hConvF, err := slide.Conv(f, h)
	if err != nil {
		t.Fatal(err)
	}
	hCorrF, err := slide.Corr(f, h)
	if err != nil {
		t.Fatal(err)
	}
	if err := errIfNotEqImage(gCorrF, hConvF, eps); err != nil {
		t.Fatal(err)
	}
	if err := errIfNotEqImage(hCorrF, gConvF, eps); err != nil {
		t.Fatal(err)
	}
}

func TestCorrFFT_vsNaive(t *testing.T) {
	const (
		m   = 40
		n   = 30
		w   = 100
		h   = 80
		eps = 1e-9
	)
	f := randImage(w, h)
	g := randImage(m, n)
	naive, err := slide.CorrNaive(f, g)
	if err != nil {
		t.Fatal(err)
	}
	fft, err := slide.CorrFFT(f, g)
	if err != nil {
		t.Fatal(err)
	}
	if !naive.Size().Eq(fft.Size()) {
		t.Fatalf("size mismatch (naive %v, fft %v)", naive.Size(), fft.Size())
	}
	if err := errIfNotEqImage(naive, fft, eps); err != nil {
		t.Fatal(err)
	}
}

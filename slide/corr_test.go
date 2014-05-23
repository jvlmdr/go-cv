package slide

import (
	"image"
	"math"
	"math/rand"
	"testing"

	"github.com/jackvalmadre/go-cv/rimg64"
)

func TestCorr(t *testing.T) {
	const eps = 1e-9

	f := rimg64.New(5, 3)
	g := rimg64.New(3, 2)
	// f:
	// 1 2 3 4 5
	// 2 5 4 1 3
	// 5 4 3 2 1
	f.Set(0, 0, 1)
	f.Set(1, 0, 2)
	f.Set(2, 0, 3)
	f.Set(3, 0, 4)
	f.Set(4, 0, 5)
	f.Set(0, 1, 2)
	f.Set(1, 1, 5)
	f.Set(2, 1, 4)
	f.Set(3, 1, 1)
	f.Set(4, 1, 3)
	f.Set(0, 2, 5)
	f.Set(1, 2, 4)
	f.Set(2, 2, 3)
	f.Set(3, 2, 2)
	f.Set(4, 2, 1)
	// g:
	// 3 1 5
	// 2 4 1
	g.Set(0, 0, 3)
	g.Set(1, 0, 1)
	g.Set(2, 0, 5)
	g.Set(0, 1, 2)
	g.Set(1, 1, 4)
	g.Set(2, 1, 1)

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

	h := Corr(f, g)
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

// Compare naive and Fourier implementations.
func TestCorr_FFTVsNaive(t *testing.T) {
	const (
		m   = 40
		n   = 30
		w   = 100
		h   = 80
		eps = 1e-12
	)

	f := randImage(w, h)
	g := randImage(m, n)

	naive := corrNaive(f, g)
	fourier := corrFFT(f, g)

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

func randImage(width, height int) *rimg64.Image {
	f := rimg64.New(width, height)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			f.Set(i, j, rand.NormFloat64())
		}
	}
	return f
}

func randMulti(width, height, channels int) *rimg64.Multi {
	f := rimg64.NewMulti(width, height, channels)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			for k := 0; k < channels; k++ {
				f.Set(i, j, k, rand.NormFloat64())
			}
		}
	}
	return f
}

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

	naive := corrMultiNaive(f, g)
	fourier := corrMultiFFT(f, g)

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

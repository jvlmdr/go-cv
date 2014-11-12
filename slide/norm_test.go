package slide_test

import (
	"math"
	"testing"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-cv/slide"
)

func sqr(x float64) float64 { return x * x }

func TestCosCorr(t *testing.T) {
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
	gnorm := math.Sqrt(sqr(3) + sqr(1) + sqr(5) + sqr(2) + sqr(4) + sqr(1))
	cases := []struct {
		I, J int
		Want float64
	}{
		// <g, [1 2 3; 2 5 4]>
		{0, 0, (3*1 + 1*2 + 5*3 + 2*2 + 4*5 + 1*4) / math.Sqrt(sqr(1)+sqr(2)+sqr(3)+sqr(2)+sqr(5)+sqr(4)) / gnorm},
		// <g, [2 3 4; 5 4 1]>
		{1, 0, (3*2 + 1*3 + 5*4 + 2*5 + 4*4 + 1*1) / math.Sqrt(sqr(2)+sqr(3)+sqr(4)+sqr(5)+sqr(4)+sqr(1)) / gnorm},
		// <g, ([3 4 5; 4 1 3]>
		{2, 0, (3*3 + 1*4 + 5*5 + 2*4 + 4*1 + 1*3) / math.Sqrt(sqr(3)+sqr(4)+sqr(5)+sqr(4)+sqr(1)+sqr(3)) / gnorm},
		// <g, ([2 5 4; 5 4 3]>
		{0, 1, (3*2 + 1*5 + 5*4 + 2*5 + 4*4 + 1*3) / math.Sqrt(sqr(2)+sqr(5)+sqr(4)+sqr(5)+sqr(4)+sqr(3)) / gnorm},
		// <g, ([5 4 1; 4 3 2]>
		{1, 1, (3*5 + 1*4 + 5*1 + 2*4 + 4*3 + 1*2) / math.Sqrt(sqr(5)+sqr(4)+sqr(1)+sqr(4)+sqr(3)+sqr(2)) / gnorm},
		// <g, ([4 1 3; 3 2 1]>
		{2, 1, (3*4 + 1*1 + 5*3 + 2*3 + 4*2 + 1*1) / math.Sqrt(sqr(4)+sqr(1)+sqr(3)+sqr(3)+sqr(2)+sqr(1)) / gnorm},
	}

	h := slide.CosCorr(f, g, slide.FFT)
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

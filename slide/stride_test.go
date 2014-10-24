package slide_test

import (
	"image"
	"math"
	"testing"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-cv/slide"
)

func TestCorrStride(t *testing.T) {
	const eps = 1e-12
	const X = 10

	cases := []struct {
		F, G, H *rimg64.Image
		K       int
	}{
		{
			F: rimg64.FromRows([][]float64{
				{1, 2, 3, 4, 5},
				{2, 5, 4, 1, 3},
				{5, 4, 3, 2, 1},
			}),
			G: rimg64.FromRows([][]float64{
				{1, -1},
			}),
			K: 2,
			H: rimg64.FromRows([][]float64{
				{1 - 2, 3 - 4},
				{5 - 4, 3 - 2},
			}),
		},
		{
			F: rimg64.FromRows([][]float64{
				{1, 2, 3, 4, 5},
				{2, 5, 4, 1, 3},
				{5, 4, 3, 2, 1},
			}),
			G: rimg64.FromRows([][]float64{
				{1},
				{-1},
			}),
			K: 2,
			H: rimg64.FromRows([][]float64{
				{1 - 2, 3 - 4, 5 - 3},
			}),
		},
		{
			F: rimg64.FromRows([][]float64{
				{1, 2, 3, 4, 5},
				{2, 5, 4, 1, 3},
				{5, 4, 3, 2, 1},
				{2, 2, 2, 2, 2},
			}),
			G: rimg64.FromRows([][]float64{
				{1, 2, 3},
				{-3, -2, -1},
			}),
			K: 2,
			H: rimg64.FromRows([][]float64{
				{(1*1 + 2*2 + 3*3) + ((-3)*2 + (-2)*5 + (-1)*4), (1*3 + 2*4 + 3*5) + ((-3)*4 + (-2)*1 + (-1)*3)},
				{(1*5 + 2*4 + 3*3) + ((-3)*2 + (-2)*2 + (-1)*2), (1*3 + 2*2 + 3*1) + ((-3)*2 + (-2)*2 + (-1)*2)},
			}),
		},
		{
			F: rimg64.FromRows([][]float64{
				{1, 2, 3, 4, 5},
				{2, 5, 4, 1, 3},
				{5, 4, 3, 2, 1},
				{2, 2, 2, 2, 2},
			}),
			G: rimg64.FromRows([][]float64{
				{1, 2, 3},
				{-3, -2, -1},
			}),
			K: 3,
			H: rimg64.FromRows([][]float64{
				{(1*1 + 2*2 + 3*3) + ((-3)*2 + (-2)*5 + (-1)*4)},
			}),
		},
		{
			F: rimg64.FromRows([][]float64{
				{1, 2, 3, 4, 5},
				{2, 5, 4, 1, 3},
				{5, 4, 3, 2, 1},
				{2, 2, 2, 2, 2},
			}),
			G: rimg64.FromRows([][]float64{
				{1, 2},
				{-3, -2},
			}),
			K: 3,
			H: rimg64.FromRows([][]float64{
				{(1*1 + 2*2) + ((-3)*2 + (-2)*5), (1*4 + 2*5) + ((-3)*1 + (-2)*3)},
			}),
		},

		{
			F: rimg64.FromRows([][]float64{
				{1, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 1, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 1, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 1},
				{0, 0, 0, 0, 0, 0, 0, 1},
			}),
			G: rimg64.FromRows([][]float64{
				{1, 2},
				{3, 4},
			}),
			K: 2,
			H: rimg64.FromRows([][]float64{
				{1, 0, 3, 0},
				{0, 0, 0, 0},
				{0, 1, 0, 4},
			}),
		},
		{
			F: rimg64.FromRows([][]float64{
				{1, 0, X, 0, 1, X, 0, 0},
				{0, 0, X, 0, 0, X, 1, 0},
				{X, X, X, X, X, X, X, X},
				{0, 0, X, 1, 0, X, 0, 1},
				{0, 1, X, 0, 1, X, 0, 1},
				{X, X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X, X},
			}),
			G: rimg64.FromRows([][]float64{
				{1, 2},
				{3, 4},
			}),
			K: 3,
			H: rimg64.FromRows([][]float64{
				{1, 2, 3},
				{4, 5, 6},
			}),
		},
		{
			F: rimg64.FromRows([][]float64{
				{0, 0, X, X, 0, 1, X, X},
				{0, 1, X, X, 0, 0, X, X},
				{X, X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X, X},
				{0, 0, X, X, 1, 0, X, X},
				{1, 0, X, X, 0, 1, X, X},
				{X, X, X, X, X, X, X, X},
			}),
			G: rimg64.FromRows([][]float64{
				{1, 2},
				{3, 4},
			}),
			K: 4,
			H: rimg64.FromRows([][]float64{
				{4, 2},
				{3, 5},
			}),
		},
		{
			F: rimg64.FromRows([][]float64{
				{0, 1, X, X, X, 0, 0, X},
				{0, 0, X, X, X, 0, 1, X},
				{X, X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X, X},
				{0, 1, X, X, X, 1, 0, X},
				{0, 1, X, X, X, 0, 0, X},
			}),
			G: rimg64.FromRows([][]float64{
				{1, 2},
				{3, 4},
			}),
			K: 5,
			H: rimg64.FromRows([][]float64{
				{2, 4},
				{6, 1},
			}),
		},
		{
			F: rimg64.FromRows([][]float64{
				{1, 1, X, X, X, X, 0, 0},
				{0, 0, X, X, X, X, 1, 1},
				{X, X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X, X},
			}),
			G: rimg64.FromRows([][]float64{
				{1, 2},
				{3, 4},
			}),
			K: 6,
			H: rimg64.FromRows([][]float64{
				{3, 7},
			}),
		},
		{
			F: rimg64.FromRows([][]float64{
				{1, 1, X, X, X, X, X, X},
				{1, 1, X, X, X, X, X, X},
				{X, X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X, X},
			}),
			G: rimg64.FromRows([][]float64{
				{1, 2},
				{3, 4},
			}),
			K: 7,
			H: rimg64.FromRows([][]float64{
				{10},
			}),
		},
		{
			F: rimg64.FromRows([][]float64{
				{1, 1, X, X, X, X, X, X},
				{0, 1, X, X, X, X, X, X},
				{X, X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X, X},
				{X, X, X, X, X, X, X, X},
			}),
			G: rimg64.FromRows([][]float64{
				{1, 2},
				{3, 4},
			}),
			K: 10000,
			H: rimg64.FromRows([][]float64{
				{7},
			}),
		},

		{
			F: rimg64.FromRows([][]float64{
				{1, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 1, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 1, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 1, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
			}),
			G: rimg64.FromRows([][]float64{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			}),
			K: 2,
			H: rimg64.FromRows([][]float64{
				{1, 6, 4},
				{9, 7, 0},
				{3, 1, 5},
			}),
		},
	}

	for _, q := range cases {
		h := slide.CorrStride(q.F, q.G, q.K)
		t.Logf("f %v, g %v, k %d", q.F.Size(), q.G.Size(), q.K)
		checkEq(t, q.H, h, eps)
	}
}

func TestCorrStride_rand(t *testing.T) {
	const eps = 1e-12
	cases := []struct {
		F, G image.Point
		K    int
	}{
		{F: image.Pt(8, 10), G: image.Pt(3, 2), K: 5},
		{F: image.Pt(100, 1), G: image.Pt(1, 1), K: 5},
		{F: image.Pt(1, 100), G: image.Pt(1, 1), K: 5},
		{F: image.Pt(43, 64), G: image.Pt(4, 5), K: 3},
		{F: image.Pt(43, 64), G: image.Pt(5, 4), K: 3},
		{F: image.Pt(64, 43), G: image.Pt(4, 5), K: 3},
		{F: image.Pt(64, 43), G: image.Pt(5, 4), K: 3},
		{F: image.Pt(63, 127), G: image.Pt(3, 2), K: 32},
		{F: image.Pt(63, 127), G: image.Pt(2, 3), K: 32},
		{F: image.Pt(63, 127), G: image.Pt(3, 2), K: 31},
		{F: image.Pt(63, 127), G: image.Pt(2, 3), K: 31},
		{F: image.Pt(63, 127), G: image.Pt(2, 3), K: 10000},
	}

	for _, q := range cases {
		f := randImage(q.F.X, q.F.Y)
		g := randImage(q.G.X, q.G.Y)
		want := slide.Decimate(slide.Corr(f, g), q.K)
		got := slide.CorrStride(f, g, q.K)
		t.Logf("f %v, g %v, k %d", q.F, q.G, q.K)
		checkEq(t, want, got, eps)
	}
}

func TestCorrMultiStride_rand(t *testing.T) {
	const eps = 1e-12
	cases := []struct {
		F, G image.Point
		C    int
		K    int
	}{
		{F: image.Pt(8, 10), G: image.Pt(3, 2), C: 5, K: 5},
		{F: image.Pt(100, 1), G: image.Pt(1, 1), C: 5, K: 5},
		{F: image.Pt(1, 100), G: image.Pt(1, 1), C: 5, K: 5},
		{F: image.Pt(43, 64), G: image.Pt(4, 5), C: 5, K: 3},
		{F: image.Pt(43, 64), G: image.Pt(5, 4), C: 5, K: 3},
		{F: image.Pt(64, 43), G: image.Pt(4, 5), C: 5, K: 3},
		{F: image.Pt(64, 43), G: image.Pt(5, 4), C: 5, K: 3},
		{F: image.Pt(63, 127), G: image.Pt(3, 2), C: 5, K: 32},
		{F: image.Pt(63, 127), G: image.Pt(2, 3), C: 5, K: 32},
		{F: image.Pt(63, 127), G: image.Pt(3, 2), C: 5, K: 31},
		{F: image.Pt(63, 127), G: image.Pt(2, 3), C: 5, K: 31},
		{F: image.Pt(63, 127), G: image.Pt(2, 3), C: 5, K: 10000},
	}

	for _, q := range cases {
		f := randMulti(q.F.X, q.F.Y, q.C)
		g := randMulti(q.G.X, q.G.Y, q.C)
		want := slide.Decimate(slide.CorrMulti(f, g), q.K)
		got := slide.CorrMultiStride(f, g, q.K)
		t.Logf("f %v, g %v, k %d", q.F, q.G, q.K)
		checkEq(t, want, got, eps)
	}
}

func checkEq(t *testing.T, want, got *rimg64.Image, eps float64) {
	if !want.Size().Eq(got.Size()) {
		t.Errorf("size: want %v, got %v", want.Size(), got.Size())
		return
	}
	for x := 0; x < want.Width; x++ {
		for y := 0; y < want.Height; y++ {
			a, b := want.At(x, y), got.At(x, y)
			if math.Abs(a-b) > eps {
				t.Errorf("value at %v: want %g, got %g", image.Pt(x, y), a, b)
			}
		}
	}
}

package slide_test

import (
	"image"
	"testing"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-cv/slide"
)

func TestCorrStrideNaive(t *testing.T) {
	const eps = 1e-9
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
		h, err := slide.CorrStrideNaive(q.F, q.G, q.K)
		if err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.F.Size(), q.G.Size(), q.K, err)
			continue
		}
		if err := errIfNotEqImage(q.H, h, eps); err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.F.Size(), q.G.Size(), q.K, err)
			continue
		}
	}
}

var strideCases = []struct {
	ImSize   image.Point
	TmplSize image.Point
	K        int
}{
	{ImSize: image.Pt(8, 10), TmplSize: image.Pt(3, 2), K: 5},
	{ImSize: image.Pt(100, 1), TmplSize: image.Pt(1, 1), K: 5},
	{ImSize: image.Pt(1, 100), TmplSize: image.Pt(1, 1), K: 5},
	{ImSize: image.Pt(43, 64), TmplSize: image.Pt(4, 5), K: 3},
	{ImSize: image.Pt(43, 64), TmplSize: image.Pt(5, 4), K: 3},
	{ImSize: image.Pt(64, 43), TmplSize: image.Pt(4, 5), K: 3},
	{ImSize: image.Pt(64, 43), TmplSize: image.Pt(5, 4), K: 3},
	{ImSize: image.Pt(63, 127), TmplSize: image.Pt(3, 2), K: 32},
	{ImSize: image.Pt(63, 127), TmplSize: image.Pt(2, 3), K: 32},
	{ImSize: image.Pt(63, 127), TmplSize: image.Pt(3, 2), K: 31},
	{ImSize: image.Pt(63, 127), TmplSize: image.Pt(2, 3), K: 31},
	{ImSize: image.Pt(63, 127), TmplSize: image.Pt(2, 3), K: 10000},
}

func TestCorrStrideNaive_vsDecimate(t *testing.T) {
	const eps = 1e-9
	for _, q := range strideCases {
		f := randImage(q.ImSize.X, q.ImSize.Y)
		g := randImage(q.TmplSize.X, q.TmplSize.Y)
		h, err := slide.CorrNaive(f, g)
		if err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.ImSize, q.TmplSize, q.K, err)
			continue
		}
		want := slide.Decimate(h, q.K)
		got, err := slide.CorrStrideNaive(f, g, q.K)
		if err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.ImSize, q.TmplSize, q.K, err)
			continue
		}
		if err := errIfNotEqImage(want, got, eps); err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.ImSize, q.TmplSize, q.K, err)
			continue
		}
	}
}

func TestCorrStrideFFT_vsNaive(t *testing.T) {
	const eps = 1e-9
	for _, q := range strideCases {
		f := randImage(q.ImSize.X, q.ImSize.Y)
		g := randImage(q.TmplSize.X, q.TmplSize.Y)
		naive, err := slide.CorrStrideNaive(f, g, q.K)
		if err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.ImSize, q.TmplSize, q.K, err)
			continue
		}
		fft, err := slide.CorrStrideFFT(f, g, q.K)
		if err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.ImSize, q.TmplSize, q.K, err)
			continue
		}
		if err := errIfNotEqImage(naive, fft, eps); err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.ImSize, q.TmplSize, q.K, err)
			continue
		}
	}
}

func TestCorrStrideBLAS_vsNaive(t *testing.T) {
	const eps = 1e-9
	for _, q := range strideCases {
		f := randImage(q.ImSize.X, q.ImSize.Y)
		g := randImage(q.TmplSize.X, q.TmplSize.Y)
		naive, err := slide.CorrStrideNaive(f, g, q.K)
		if err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.ImSize, q.TmplSize, q.K, err)
			continue
		}
		blas, err := slide.CorrStrideBLAS(f, g, q.K)
		if err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.ImSize, q.TmplSize, q.K, err)
			continue
		}
		if err := errIfNotEqImage(naive, blas, eps); err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.ImSize, q.TmplSize, q.K, err)
			continue
		}
	}
}

func BenchmarkCorrStrideNaive_Im_640x480_Tmpl_3x3_Stride_4(b *testing.B) {
	benchmarkCorrStride(b, image.Pt(640, 480), image.Pt(3, 3), 4, slide.Naive)
}

func BenchmarkCorrStrideNaive_Im_640x480_Tmpl_16x16_Stride_4(b *testing.B) {
	benchmarkCorrStride(b, image.Pt(640, 480), image.Pt(16, 16), 4, slide.Naive)
}

func BenchmarkCorrStrideFFT_Im_640x480_Tmpl_3x3_Stride_4(b *testing.B) {
	benchmarkCorrStride(b, image.Pt(640, 480), image.Pt(3, 3), 4, slide.FFT)
}

func BenchmarkCorrStrideFFT_Im_640x480_Tmpl_16x16_Stride_4(b *testing.B) {
	benchmarkCorrStride(b, image.Pt(640, 480), image.Pt(16, 16), 4, slide.FFT)
}

func BenchmarkCorrStrideBLAS_Im_640x480_Tmpl_3x3_Stride_4(b *testing.B) {
	benchmarkCorrStride(b, image.Pt(640, 480), image.Pt(3, 3), 4, slide.BLAS)
}

func BenchmarkCorrStrideBLAS_Im_640x480_Tmpl_16x16_Stride_4(b *testing.B) {
	benchmarkCorrStride(b, image.Pt(640, 480), image.Pt(16, 16), 4, slide.BLAS)
}

func benchmarkCorrStride(b *testing.B, im, tmpl image.Point, stride int, algo slide.Algo) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		f := randImage(im.X, im.Y)
		g := randImage(tmpl.X, tmpl.Y)
		b.StartTimer()
		slide.CorrStrideAlgo(f, g, stride, algo)
	}
}

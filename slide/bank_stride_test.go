package slide_test

import (
	"image"
	"testing"

	"github.com/jvlmdr/go-cv/slide"
)

func TestCorrBankStrideNaive_vsDecimate(t *testing.T) {
	const (
		numOut = 4
		eps    = 1e-9
	)
	for _, q := range strideCases {
		f := randImage(q.ImSize.X, q.ImSize.Y)
		g := randBank(q.TmplSize.X, q.TmplSize.Y, numOut)
		h, err := slide.CorrBankNaive(f, g)
		if err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.ImSize, q.TmplSize, q.K, err)
			continue
		}
		want := slide.DecimateMulti(h, q.K)
		got, err := slide.CorrBankStrideNaive(f, g, q.K)
		if err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.ImSize, q.TmplSize, q.K, err)
			continue
		}
		if err := errIfNotEqMulti(want, got, eps); err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.ImSize, q.TmplSize, q.K, err)
			continue
		}
	}
}

func TestCorrBankStrideFFT_vsNaive(t *testing.T) {
	const (
		m      = 40
		n      = 30
		w      = 100
		h      = 80
		numOut = 6
		stride = 3
		eps    = 1e-9
	)
	f := randImage(w, h)
	g := randBank(w, h, numOut)
	naive, err := slide.CorrBankStrideNaive(f, g, stride)
	if err != nil {
		t.Fatal(err)
	}
	fft, err := slide.CorrBankStrideFFT(f, g, stride)
	if err != nil {
		t.Fatal(err)
	}
	if err := errIfNotEqMulti(naive, fft, eps); err != nil {
		t.Fatal(err)
	}
}

func TestCorrBankStrideBLAS_vsNaive(t *testing.T) {
	const (
		m      = 40
		n      = 30
		w      = 100
		h      = 80
		numOut = 6
		stride = 3
		eps    = 1e-9
	)
	f := randImage(w, h)
	g := randBank(w, h, numOut)
	naive, err := slide.CorrBankStrideNaive(f, g, stride)
	if err != nil {
		t.Fatal(err)
	}
	blas, err := slide.CorrBankStrideBLAS(f, g, stride)
	if err != nil {
		t.Fatal(err)
	}
	if err := errIfNotEqMulti(naive, blas, eps); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkCorrBankStrideNaive_Im_640x480_Tmpl_3x3_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 4, 4, slide.Naive)
}

func BenchmarkCorrBankStrideNaive_Im_640x480_Tmpl_3x3_Out_32_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 32, 4, slide.Naive)
}

func BenchmarkCorrBankStrideNaive_Im_640x480_Tmpl_16x16_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(16, 16), 4, 4, slide.Naive)
}

func BenchmarkCorrBankStrideFFT_Im_640x480_Tmpl_3x3_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 4, 4, slide.FFT)
}

func BenchmarkCorrBankStrideFFT_Im_640x480_Tmpl_3x3_Out_32_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 32, 4, slide.FFT)
}

func BenchmarkCorrBankStrideFFT_Im_640x480_Tmpl_16x16_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(16, 16), 4, 4, slide.FFT)
}

func BenchmarkCorrBankStrideFFT_Im_640x480_Tmpl_16x16_Out_96_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(16, 16), 96, 4, slide.FFT)
}

func BenchmarkCorrBankStrideBLAS_Im_640x480_Tmpl_3x3_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 4, 4, slide.BLAS)
}

func BenchmarkCorrBankStrideBLAS_Im_640x480_Tmpl_3x3_Out_32_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 32, 4, slide.BLAS)
}

func BenchmarkCorrBankStrideBLAS_Im_640x480_Tmpl_16x16_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(16, 16), 4, 4, slide.BLAS)
}

func BenchmarkCorrBankStrideBLAS_Im_640x480_Tmpl_16x16_Out_96_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(16, 16), 96, 4, slide.BLAS)
}

func benchmarkCorrBankStride(b *testing.B, im, tmpl image.Point, out, stride int, algo slide.Algo) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		f := randImage(im.X, im.Y)
		g := randBank(tmpl.X, tmpl.Y, out)
		b.StartTimer()
		slide.CorrBankStrideAlgo(f, g, stride, algo)
	}
}

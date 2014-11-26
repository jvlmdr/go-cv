package slide_test

import (
	"image"
	"testing"

	"github.com/jvlmdr/go-cv/slide"
)

func TestCorrBankStrideFFT_vsNaive(t *testing.T) {
	const (
		m      = 40
		n      = 30
		w      = 100
		h      = 80
		numOut = 6
		r      = 3
		eps    = 1e-9
	)
	f := randImage(w, h)
	g := randBank(w, h, numOut)
	naive, err := slide.CorrBankStrideNaive(f, g, r)
	if err != nil {
		t.Fatal(err)
	}
	fft, err := slide.CorrBankStrideFFT(f, g, r)
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
		r      = 3
		eps    = 1e-9
	)
	f := randImage(w, h)
	g := randBank(w, h, numOut)
	naive, err := slide.CorrBankStrideNaive(f, g, r)
	if err != nil {
		t.Fatal(err)
	}
	blas, err := slide.CorrBankStrideBLAS(f, g, r)
	if err != nil {
		t.Fatal(err)
	}
	if err := errIfNotEqMulti(naive, blas, eps); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkCorrBankStrideFFT_Im_640x480_Tmpl_3x3_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 4, 4, true, false)
}

func BenchmarkCorrBankStrideFFT_Im_640x480_Tmpl_3x3_Out_32_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 32, 4, true, false)
}

func BenchmarkCorrBankStrideFFT_Im_640x480_Tmpl_16x16_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(16, 16), 4, 4, true, false)
}

func BenchmarkCorrBankStrideBLAS_Im_640x480_Tmpl_3x3_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 4, 4, false, true)
}

func BenchmarkCorrBankStrideBLAS_Im_640x480_Tmpl_3x3_Out_32_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 32, 4, false, true)
}

func BenchmarkCorrBankStrideBLAS_Im_640x480_Tmpl_16x16_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(16, 16), 4, 4, false, true)
}

func BenchmarkCorrBankStrideNaive_Im_640x480_Tmpl_3x3_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 4, 4, false, false)
}

func BenchmarkCorrBankStrideNaive_Im_640x480_Tmpl_3x3_Out_32_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 32, 4, false, false)
}

func BenchmarkCorrBankStrideNaive_Im_640x480_Tmpl_16x16_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(16, 16), 4, 4, false, false)
}

func BenchmarkCorrBankStrideFFT_Im_640x480_Tmpl_16x16_Out_96_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(16, 16), 96, 4, true, false)
}

func BenchmarkCorrBankStrideBLAS_Im_640x480_Tmpl_16x16_Out_96_Stride_4(b *testing.B) {
	benchmarkCorrBankStride(b, image.Pt(640, 480), image.Pt(16, 16), 96, 4, false, true)
}

func benchmarkCorrBankStride(b *testing.B, im, tmpl image.Point, out, stride int, fft, blas bool) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		f := randImage(im.X, im.Y)
		g := randBank(tmpl.X, tmpl.Y, out)
		b.StartTimer()
		if fft {
			slide.CorrBankStrideFFT(f, g, stride)
		} else if blas {
			slide.CorrBankStrideBLAS(f, g, stride)
		} else {
			slide.CorrBankStrideNaive(f, g, stride)
		}
	}
}

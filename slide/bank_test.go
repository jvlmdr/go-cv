package slide_test

import (
	"image"
	"testing"

	"github.com/jvlmdr/go-cv/slide"
)

func TestCorrBankFFT_vsNaive(t *testing.T) {
	const (
		m      = 40
		n      = 30
		w      = 100
		h      = 80
		numOut = 6
		eps    = 1e-9
	)
	f := randImage(w, h)
	g := randBank(w, h, numOut)
	naive, err := slide.CorrBankNaive(f, g)
	if err != nil {
		t.Fatal(err)
	}
	fft, err := slide.CorrBankFFT(f, g)
	if err != nil {
		t.Fatal(err)
	}
	if err := errIfNotEqMulti(naive, fft, eps); err != nil {
		t.Fatal(err)
	}
}

func TestCorrBankBLAS_vsNaive(t *testing.T) {
	const (
		m      = 40
		n      = 30
		w      = 100
		h      = 80
		numOut = 6
		eps    = 1e-9
	)
	f := randImage(w, h)
	g := randBank(w, h, numOut)
	naive, err := slide.CorrBankNaive(f, g)
	if err != nil {
		t.Fatal(err)
	}
	blas, err := slide.CorrBankBLAS(f, g)
	if err != nil {
		t.Fatal(err)
	}
	if err := errIfNotEqMulti(naive, blas, eps); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkCorrBankNaive_Im_640x480_Tmpl_3x3_Out_4(b *testing.B) {
	benchmarkCorrBank(b, image.Pt(640, 480), image.Pt(3, 3), 4, slide.Naive)
}

func BenchmarkCorrBankNaive_Im_640x480_Tmpl_3x3_Out_32(b *testing.B) {
	benchmarkCorrBank(b, image.Pt(640, 480), image.Pt(3, 3), 32, slide.Naive)
}

func BenchmarkCorrBankNaive_Im_640x480_Tmpl_16x16_Out_4(b *testing.B) {
	benchmarkCorrBank(b, image.Pt(640, 480), image.Pt(16, 16), 4, slide.Naive)
}

func BenchmarkCorrBankFFT_Im_640x480_Tmpl_3x3_Out_4(b *testing.B) {
	benchmarkCorrBank(b, image.Pt(640, 480), image.Pt(3, 3), 4, slide.FFT)
}

func BenchmarkCorrBankFFT_Im_640x480_Tmpl_3x3_Out_32(b *testing.B) {
	benchmarkCorrBank(b, image.Pt(640, 480), image.Pt(3, 3), 32, slide.FFT)
}

func BenchmarkCorrBankFFT_Im_640x480_Tmpl_16x16_Out_4(b *testing.B) {
	benchmarkCorrBank(b, image.Pt(640, 480), image.Pt(16, 16), 4, slide.FFT)
}

func BenchmarkCorrBankBLAS_Im_640x480_Tmpl_3x3_Out_4(b *testing.B) {
	benchmarkCorrBank(b, image.Pt(640, 480), image.Pt(3, 3), 4, slide.BLAS)
}

func BenchmarkCorrBankBLAS_Im_640x480_Tmpl_3x3_Out_32(b *testing.B) {
	benchmarkCorrBank(b, image.Pt(640, 480), image.Pt(3, 3), 32, slide.BLAS)
}

func BenchmarkCorrBankBLAS_Im_640x480_Tmpl_16x16_Out_4(b *testing.B) {
	benchmarkCorrBank(b, image.Pt(640, 480), image.Pt(16, 16), 4, slide.BLAS)
}

func benchmarkCorrBank(b *testing.B, im, tmpl image.Point, out int, algo slide.Algo) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		f := randImage(im.X, im.Y)
		g := randBank(tmpl.X, tmpl.Y, out)
		b.StartTimer()
		slide.CorrBankAlgo(f, g, algo)
	}
}

package slide_test

import (
	"image"
	"testing"

	"github.com/jvlmdr/go-cv/slide"
)

func TestCorrMultiFFT_vsNaive(t *testing.T) {
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
	naive, err := slide.CorrMultiNaive(f, g)
	if err != nil {
		t.Fatal(err)
	}
	fft, err := slide.CorrMultiFFT(f, g)
	if err != nil {
		t.Fatal(err)
	}
	if err := errIfNotEqImage(naive, fft, eps); err != nil {
		t.Fatal(err)
	}
}

func TestCorrMultiBLAS_vsNaive(t *testing.T) {
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
	naive, err := slide.CorrMultiNaive(f, g)
	if err != nil {
		t.Fatal(err)
	}
	blas, err := slide.CorrMultiBLAS(f, g)
	if err != nil {
		t.Fatal(err)
	}
	if err := errIfNotEqImage(naive, blas, eps); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkCorrMultiNaive_Im_640x480_Tmpl_3x3_In_4(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(3, 3), 4, slide.Naive)
}

func BenchmarkCorrMultiNaive_Im_640x480_Tmpl_3x3_In_32(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(3, 3), 32, slide.Naive)
}

func BenchmarkCorrMultiNaive_Im_640x480_Tmpl_16x16_In_4(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(16, 16), 4, slide.Naive)
}

func BenchmarkCorrMultiFFT_Im_640x480_Tmpl_3x3_In_4(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(3, 3), 4, slide.FFT)
}

func BenchmarkCorrMultiFFT_Im_640x480_Tmpl_3x3_In_32(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(3, 3), 32, slide.FFT)
}

func BenchmarkCorrMultiFFT_Im_640x480_Tmpl_16x16_In_4(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(16, 16), 4, slide.FFT)
}

func BenchmarkCorrMultiFFT_Im_640x480_Tmpl_16x16_In_32(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(16, 16), 32, slide.FFT)
}

func BenchmarkCorrMultiBLAS_Im_640x480_Tmpl_3x3_In_4(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(3, 3), 4, slide.BLAS)
}

func BenchmarkCorrMultiBLAS_Im_640x480_Tmpl_3x3_In_32(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(3, 3), 32, slide.BLAS)
}

func BenchmarkCorrMultiBLAS_Im_640x480_Tmpl_16x16_In_4(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(16, 16), 4, slide.BLAS)
}

func BenchmarkCorrMultiBLAS_Im_640x480_Tmpl_16x16_In_32(b *testing.B) {
	if testing.Short() {
		b.Skip("skip: 16x16 template, 32 input channels")
	}
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(16, 16), 32, slide.BLAS)
}

func benchmarkCorrMulti(b *testing.B, im, tmpl image.Point, c int, algo slide.Algo) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		f := randMulti(im.X, im.Y, c)
		g := randMulti(tmpl.X, tmpl.Y, c)
		b.StartTimer()
		slide.CorrMultiAlgo(f, g, algo)
	}
}

package slide_test

import (
	"image"
	"testing"

	"github.com/jvlmdr/go-cv/slide"
)

// Compare naive and Fourier implementations.
func TestCorrMultiBankFFT_vsNaive(t *testing.T) {
	const (
		m   = 40
		n   = 30
		w   = 100
		h   = 80
		p   = 4
		q   = 6
		eps = 1e-12
	)

	f := randMulti(w, h, p)
	g := randMultiBank(w, h, p, q)

	naive, err := slide.CorrMultiBankNaive(f, g)
	if err != nil {
		t.Fatal(err)
	}
	fourier, err := slide.CorrMultiBankFFT(f, g)
	if err != nil {
		t.Fatal(err)
	}
	if err := errIfNotEqMulti(naive, fourier, 1e-9); err != nil {
		t.Fatal(err)
	}
}

// Compare naive and Fourier implementations.
func TestCorrMultiBankBLAS_vsNaive(t *testing.T) {
	const (
		m   = 40
		n   = 30
		w   = 100
		h   = 80
		p   = 4
		q   = 6
		eps = 1e-12
	)

	f := randMulti(w, h, p)
	g := randMultiBank(w, h, p, q)

	naive, err := slide.CorrMultiBankNaive(f, g)
	if err != nil {
		t.Fatal(err)
	}
	blas, err := slide.CorrMultiBankBLAS(f, g)
	if err != nil {
		t.Fatal(err)
	}
	if err := errIfNotEqMulti(naive, blas, 1e-9); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkCorrMultiBankFFT_640x480_3x3_4_4(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 4, 4, true, false)
}

func BenchmarkCorrMultiBankFFT_640x480_3x3_4_128(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 4, 128, true, false)
}

func BenchmarkCorrMultiBankFFT_640x480_3x3_128_4(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 128, 4, true, false)
}

func BenchmarkCorrMultiBankFFT_640x480_16x16_4_4(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(16, 16), 4, 4, true, false)
}

func BenchmarkCorrMultiBankBLAS_640x480_3x3_4_4(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 4, 4, false, true)
}

func BenchmarkCorrMultiBankBLAS_640x480_3x3_4_128(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 4, 128, false, true)
}

func BenchmarkCorrMultiBankBLAS_640x480_3x3_128_4(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 128, 4, false, true)
}

func BenchmarkCorrMultiBankBLAS_640x480_3x3_128_128(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 128, 128, false, true)
}

func BenchmarkCorrMultiBankBLAS_640x480_16x16_4_4(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(16, 16), 4, 4, false, true)
}

func BenchmarkCorrMultiBankNaive_640x480_3x3_4_4(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 4, 4, false, false)
}

func BenchmarkCorrMultiBankNaive_640x480_3x3_4_128(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 4, 128, false, false)
}

func BenchmarkCorrMultiBankNaive_640x480_3x3_128_4(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 128, 4, false, false)
}

func BenchmarkCorrMultiBankNaive_640x480_16x16_4_4(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(16, 16), 4, 4, false, false)
}

func benchmarkCorrMultiBank(b *testing.B, im, tmpl image.Point, in, out int, fft, blas bool) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		f := randMulti(im.X, im.Y, in)
		g := randMultiBank(tmpl.X, tmpl.Y, in, out)
		b.StartTimer()
		if fft {
			slide.CorrMultiBankFFT(f, g)
		} else if blas {
			slide.CorrMultiBankBLAS(f, g)
		} else {
			slide.CorrMultiBankNaive(f, g)
		}
	}
}

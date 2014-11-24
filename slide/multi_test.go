package slide_test

import (
	"image"
	"testing"

	"github.com/jvlmdr/go-cv/slide"
)

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

func BenchmarkCorrMultiFFT_Im_640x480_Tmpl_3x3_In_4(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(3, 3), 4, true)
}

func BenchmarkCorrMultiFFT_Im_640x480_Tmpl_3x3_In_128(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(3, 3), 128, true)
}

func BenchmarkCorrMultiFFT_Im_640x480_Tmpl_16x16_In_4(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(16, 16), 4, true)
}

func BenchmarkCorrMultiFFT_Im_640x480_Tmpl_16x16_In_128(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(16, 16), 128, true)
}

func BenchmarkCorrMultiNaive_Im_640x480_Tmpl_3x3_In_4(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(3, 3), 4, false)
}

func BenchmarkCorrMultiNaive_Im_640x480_Tmpl_3x3_In_128(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(3, 3), 128, false)
}

func BenchmarkCorrMultiNaive_Im_640x480_Tmpl_16x16_In_4(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(16, 16), 4, false)
}

func benchmarkCorrMulti(b *testing.B, im, tmpl image.Point, c int, fft bool) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		f := randMulti(im.X, im.Y, c)
		g := randMulti(tmpl.X, tmpl.Y, c)
		b.StartTimer()
		if fft {
			slide.CorrMultiFFT(f, g)
		} else {
			slide.CorrMultiNaive(f, g)
		}
	}
}

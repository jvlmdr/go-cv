package slide_test

import (
	"image"
	"math"
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

	naive := slide.CorrMultiNaive(f, g)
	fourier := slide.CorrMultiFFT(f, g)

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

func BenchmarkCorrMultiFFT_640x480_3x3_4(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(3, 3), 4, true)
}

func BenchmarkCorrMultiFFT_640x480_3x3_128(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(3, 3), 128, true)
}

func BenchmarkCorrMultiFFT_640x480_16x16_4(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(16, 16), 4, true)
}

func BenchmarkCorrMultiFFT_640x480_16x16_128(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(16, 16), 128, true)
}

func BenchmarkCorrMultiNaive_640x480_3x3_4(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(3, 3), 4, false)
}

func BenchmarkCorrMultiNaive_640x480_3x3_128(b *testing.B) {
	benchmarkCorrMulti(b, image.Pt(640, 480), image.Pt(3, 3), 128, false)
}

func BenchmarkCorrMultiNaive_640x480_16x16_4(b *testing.B) {
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

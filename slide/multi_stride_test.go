package slide_test

import (
	"image"
	"testing"

	"github.com/jvlmdr/go-cv/slide"
)

func TestCorrMultiStrideNaive_vsDecimate(t *testing.T) {
	const (
		numIn = 4
		eps   = 1e-9
	)
	for _, q := range strideCases {
		f := randMulti(q.ImSize.X, q.ImSize.Y, numIn)
		g := randMulti(q.TmplSize.X, q.TmplSize.Y, numIn)
		h, err := slide.CorrMulti(f, g)
		if err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.ImSize, q.TmplSize, q.K, err)
			continue
		}
		want := slide.Decimate(h, q.K)
		got, err := slide.CorrMultiStrideNaive(f, g, q.K)
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

func TestCorrMultiStrideFFT_vsNaive(t *testing.T) {
	const (
		numIn = 4
		eps   = 1e-9
	)
	for _, q := range strideCases {
		f := randMulti(q.ImSize.X, q.ImSize.Y, numIn)
		g := randMulti(q.TmplSize.X, q.TmplSize.Y, numIn)
		naive, err := slide.CorrMultiStrideNaive(f, g, q.K)
		if err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.ImSize, q.TmplSize, q.K, err)
			continue
		}
		fft, err := slide.CorrMultiStrideFFT(f, g, q.K)
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

func TestCorrMultiStrideBLAS_vsNaive(t *testing.T) {
	const (
		numIn = 4
		eps   = 1e-9
	)
	for _, q := range strideCases {
		f := randMulti(q.ImSize.X, q.ImSize.Y, numIn)
		g := randMulti(q.TmplSize.X, q.TmplSize.Y, numIn)
		naive, err := slide.CorrMultiStrideNaive(f, g, q.K)
		if err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.ImSize, q.TmplSize, q.K, err)
			continue
		}
		blas, err := slide.CorrMultiStrideBLAS(f, g, q.K)
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

func BenchmarkCorrMultiStrideNaive_Im_640x480_Tmpl_3x3_In_4_Stride_4(b *testing.B) {
	benchmarkCorrMultiStride(b, image.Pt(640, 480), image.Pt(3, 3), 4, 4, slide.Naive)
}

func BenchmarkCorrMultiStrideNaive_Im_640x480_Tmpl_3x3_In_32_Stride_4(b *testing.B) {
	benchmarkCorrMultiStride(b, image.Pt(640, 480), image.Pt(3, 3), 32, 4, slide.Naive)
}

func BenchmarkCorrMultiStrideNaive_Im_640x480_Tmpl_16x16_In_4_Stride_4(b *testing.B) {
	benchmarkCorrMultiStride(b, image.Pt(640, 480), image.Pt(16, 16), 4, 4, slide.Naive)
}

func BenchmarkCorrMultiStrideFFT_Im_640x480_Tmpl_3x3_In_4_Stride_4(b *testing.B) {
	benchmarkCorrMultiStride(b, image.Pt(640, 480), image.Pt(3, 3), 4, 4, slide.FFT)
}

func BenchmarkCorrMultiStrideFFT_Im_640x480_Tmpl_3x3_In_32_Stride_4(b *testing.B) {
	benchmarkCorrMultiStride(b, image.Pt(640, 480), image.Pt(3, 3), 32, 4, slide.FFT)
}

func BenchmarkCorrMultiStrideFFT_Im_640x480_Tmpl_16x16_In_4_Stride_4(b *testing.B) {
	benchmarkCorrMultiStride(b, image.Pt(640, 480), image.Pt(16, 16), 4, 4, slide.FFT)
}

func BenchmarkCorrMultiStrideBLAS_Im_640x480_Tmpl_3x3_In_4_Stride_4(b *testing.B) {
	benchmarkCorrMultiStride(b, image.Pt(640, 480), image.Pt(3, 3), 4, 4, slide.BLAS)
}

func BenchmarkCorrMultiStrideBLAS_Im_640x480_Tmpl_3x3_In_32_Stride_4(b *testing.B) {
	benchmarkCorrMultiStride(b, image.Pt(640, 480), image.Pt(3, 3), 32, 4, slide.BLAS)
}

func BenchmarkCorrMultiStrideBLAS_Im_640x480_Tmpl_16x16_In_4_Stride_4(b *testing.B) {
	benchmarkCorrMultiStride(b, image.Pt(640, 480), image.Pt(16, 16), 4, 4, slide.BLAS)
}

func benchmarkCorrMultiStride(b *testing.B, im, tmpl image.Point, in, stride int, algo slide.Algo) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		f := randMulti(im.X, im.Y, in)
		g := randMulti(tmpl.X, tmpl.Y, in)
		b.StartTimer()
		slide.CorrMultiStrideAlgo(f, g, stride, algo)
	}
}

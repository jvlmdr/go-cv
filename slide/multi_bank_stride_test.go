package slide_test

import (
	"image"
	"testing"

	"github.com/jvlmdr/go-cv/slide"
)

func TestCorrMultiBankStrideFFT_vsNaive(t *testing.T) {
	const (
		m      = 40
		n      = 30
		w      = 100
		h      = 80
		numIn  = 4
		numOut = 6
		r      = 3
		eps    = 1e-9
	)
	f := randMulti(w, h, numIn)
	g := randMultiBank(w, h, numIn, numOut)
	naive, err := slide.CorrMultiBankStrideNaive(f, g, r)
	if err != nil {
		t.Fatal(err)
	}
	fft, err := slide.CorrMultiBankStrideFFT(f, g, r)
	if err != nil {
		t.Fatal(err)
	}
	if err := errIfNotEqMulti(naive, fft, eps); err != nil {
		t.Fatal(err)
	}
}

func TestCorrMultiBankStrideBLAS_vsNaive(t *testing.T) {
	const (
		m      = 40
		n      = 30
		w      = 100
		h      = 80
		numIn  = 4
		numOut = 6
		r      = 3
		eps    = 1e-9
	)
	f := randMulti(w, h, numIn)
	g := randMultiBank(w, h, numIn, numOut)
	naive, err := slide.CorrMultiBankStrideNaive(f, g, r)
	if err != nil {
		t.Fatal(err)
	}
	blas, err := slide.CorrMultiBankStrideBLAS(f, g, r)
	if err != nil {
		t.Fatal(err)
	}
	if err := errIfNotEqMulti(naive, blas, eps); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkCorrMultiBankStrideFFT_Im_640x480_Tmpl_3x3_In_4_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrMultiBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 4, 4, 4, true, false)
}

func BenchmarkCorrMultiBankStrideFFT_Im_640x480_Tmpl_3x3_In_4_Out_32_Stride_4(b *testing.B) {
	benchmarkCorrMultiBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 4, 32, 4, true, false)
}

func BenchmarkCorrMultiBankStrideFFT_Im_640x480_Tmpl_3x3_In_32_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrMultiBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 32, 4, 4, true, false)
}

func BenchmarkCorrMultiBankStrideFFT_Im_640x480_Tmpl_16x16_In_4_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrMultiBankStride(b, image.Pt(640, 480), image.Pt(16, 16), 4, 4, 4, true, false)
}

func BenchmarkCorrMultiBankStrideBLAS_Im_640x480_Tmpl_3x3_In_4_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrMultiBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 4, 4, 4, false, true)
}

func BenchmarkCorrMultiBankStrideBLAS_Im_640x480_Tmpl_3x3_In_4_Out_32_Stride_4(b *testing.B) {
	benchmarkCorrMultiBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 4, 32, 4, false, true)
}

func BenchmarkCorrMultiBankStrideBLAS_Im_640x480_Tmpl_3x3_In_32_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrMultiBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 32, 4, 4, false, true)
}

func BenchmarkCorrMultiBankStrideBLAS_Im_640x480_Tmpl_3x3_In_32_Out_32_Stride_4(b *testing.B) {
	benchmarkCorrMultiBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 32, 32, 4, false, true)
}

func BenchmarkCorrMultiBankStrideBLAS_Im_640x480_Tmpl_16x16_In_4_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrMultiBankStride(b, image.Pt(640, 480), image.Pt(16, 16), 4, 4, 4, false, true)
}

func BenchmarkCorrMultiBankStrideNaive_Im_640x480_Tmpl_3x3_In_4_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrMultiBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 4, 4, 4, false, false)
}

func BenchmarkCorrMultiBankStrideNaive_Im_640x480_Tmpl_3x3_In_4_Out_32_Stride_4(b *testing.B) {
	benchmarkCorrMultiBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 4, 32, 4, false, false)
}

func BenchmarkCorrMultiBankStrideNaive_Im_640x480_Tmpl_3x3_In_32_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrMultiBankStride(b, image.Pt(640, 480), image.Pt(3, 3), 32, 4, 4, false, false)
}

func BenchmarkCorrMultiBankStrideNaive_Im_640x480_Tmpl_16x16_In_4_Out_4_Stride_4(b *testing.B) {
	benchmarkCorrMultiBankStride(b, image.Pt(640, 480), image.Pt(16, 16), 4, 4, 4, false, false)
}

func BenchmarkCorrMultiBankStrideFFT_Im_640x480_Tmpl_16x16_In_4_Out_96_Stride_4(b *testing.B) {
	benchmarkCorrMultiBankStride(b, image.Pt(640, 480), image.Pt(16, 16), 4, 96, 4, true, false)
}

func BenchmarkCorrMultiBankStrideBLAS_Im_640x480_Tmpl_16x16_In_4_Out_96_Stride_4(b *testing.B) {
	benchmarkCorrMultiBankStride(b, image.Pt(640, 480), image.Pt(16, 16), 4, 96, 4, false, true)
}

func benchmarkCorrMultiBankStride(b *testing.B, im, tmpl image.Point, in, out, stride int, fft, blas bool) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		f := randMulti(im.X, im.Y, in)
		g := randMultiBank(tmpl.X, tmpl.Y, in, out)
		b.StartTimer()
		if fft {
			slide.CorrMultiBankStrideFFT(f, g, stride)
		} else if blas {
			slide.CorrMultiBankStrideBLAS(f, g, stride)
		} else {
			slide.CorrMultiBankStrideNaive(f, g, stride)
		}
	}
}

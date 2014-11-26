package slide_test

import (
	"image"
	"testing"

	"github.com/jvlmdr/go-cv/slide"
)

func TestCorrMultiBankFFT_vsNaive(t *testing.T) {
	const (
		m   = 40
		n   = 30
		w   = 100
		h   = 80
		p   = 4
		q   = 6
		eps = 1e-9
	)
	f := randMulti(w, h, p)
	g := randMultiBank(w, h, p, q)
	naive, err := slide.CorrMultiBankNaive(f, g)
	if err != nil {
		t.Fatal(err)
	}
	fft, err := slide.CorrMultiBankFFT(f, g)
	if err != nil {
		t.Fatal(err)
	}
	if err := errIfNotEqMulti(naive, fft, eps); err != nil {
		t.Fatal(err)
	}
}

func TestCorrMultiBankBLAS_vsNaive(t *testing.T) {
	const (
		m   = 40
		n   = 30
		w   = 100
		h   = 80
		p   = 4
		q   = 6
		eps = 1e-9
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
	if err := errIfNotEqMulti(naive, blas, eps); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkCorrMultiBankFFT_Im_640x480_Tmpl_3x3_In_4_Out_4(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 4, 4, true, false)
}

func BenchmarkCorrMultiBankFFT_Im_640x480_Tmpl_3x3_In_4_Out_32(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 4, 32, true, false)
}

func BenchmarkCorrMultiBankFFT_Im_640x480_Tmpl_3x3_In_32_Out_4(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 32, 4, true, false)
}

func BenchmarkCorrMultiBankFFT_Im_640x480_Tmpl_16x16_In_4_Out_4(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(16, 16), 4, 4, true, false)
}

func BenchmarkCorrMultiBankBLAS_Im_640x480_Tmpl_3x3_In_4_Out_4(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 4, 4, false, true)
}

func BenchmarkCorrMultiBankBLAS_Im_640x480_Tmpl_3x3_In_4_Out_32(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 4, 32, false, true)
}

func BenchmarkCorrMultiBankBLAS_Im_640x480_Tmpl_3x3_In_32_Out_4(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 32, 4, false, true)
}

func BenchmarkCorrMultiBankBLAS_Im_640x480_Tmpl_3x3_In_32_Out_32(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 32, 32, false, true)
}

func BenchmarkCorrMultiBankBLAS_Im_640x480_Tmpl_16x16_In_4_Out_4(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(16, 16), 4, 4, false, true)
}

func BenchmarkCorrMultiBankNaive_Im_640x480_Tmpl_3x3_In_4_Out_4(b *testing.B) {
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 4, 4, false, false)
}

func BenchmarkCorrMultiBankNaive_Im_640x480_Tmpl_3x3_In_4_Out_32(b *testing.B) {
	if testing.Short() {
		b.Skip("skip: 32 output channels")
	}
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 4, 32, false, false)
}

func BenchmarkCorrMultiBankNaive_Im_640x480_Tmpl_3x3_In_32_Out_4(b *testing.B) {
	if testing.Short() {
		b.Skip("skip: 32 input channels")
	}
	benchmarkCorrMultiBank(b, image.Pt(640, 480), image.Pt(3, 3), 32, 4, false, false)
}

func BenchmarkCorrMultiBankNaive_Im_640x480_Tmpl_16x16_In_4_Out_4(b *testing.B) {
	if testing.Short() {
		b.Skip("skip: 16x16 filter")
	}
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

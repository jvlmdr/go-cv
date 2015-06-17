package hog

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"testing"

	"github.com/jvlmdr/go-cv/rimg64"
)

func TestHOG_VersusFGMR(t *testing.T) {
	const (
		sbin  = 8
		fname = "000084.jpg"
	)

	// Load image.
	file, err := os.Open(fname)
	if err != nil {
		t.Fatal(err)
	}
	im, _, err := image.Decode(file)
	if err != nil {
		t.Fatal(err)
	}

	// Remove outside pixel before running C code.
	rect := im.Bounds().Inset(1).Inset(sbin / 2)
	inside := image.NewRGBA(image.Rectangle{image.ZP, rect.Size()})
	draw.Draw(inside, inside.Bounds(), im, rect.Min, draw.Src)
	// Compute transforms.
	ref := fgmr(rimg64.FromColor(inside), sbin)
	f := HOG(rimg64.FromColor(im), FGMRConfig(sbin))

	const prec = 1e-5
	// Skip first and last element. (Not using cell outside image.)
	for x := 1; x < f.Width-1; x++ {
		for y := 1; y < f.Height-1; y++ {
			for d := 0; d < f.Channels; d++ {
				want := ref.At(x, y, d)
				got := f.At(x, y, d)
				if math.Abs(want-got) > prec {
					t.Errorf("wrong value at (%d, %d, %d): want %g, got %g", x, y, d, want, got)
				}
			}
		}
	}
}

func TestHOG_boundary(t *testing.T) {
	const (
		sbin  = 4
		frac  = 3
		fname = "000084.jpg"
	)
	// Load image.
	file, err := os.Open(fname)
	if err != nil {
		t.Fatal(err)
	}
	im, _, err := image.Decode(file)

	// Take top-left part which is divisible by sbin and frac.
	im = ensureDivis(im, sbin*frac)
	// Make a rectangle of the top-left part of the image.
	// Not the most top-left window but the second-most.
	size := im.Bounds().Size()
	rect := image.Rectangle{size.Div(frac), size.Div(frac).Mul(2)}

	// Sub-sample image
	subim := image.NewRGBA(image.Rectangle{image.ZP, rect.Size()})
	draw.Draw(subim, subim.Bounds(), im, rect.Min, draw.Src)

	// Convert to real values.
	f := HOG(rimg64.FromColor(im), FGMRConfig(sbin))
	g := HOG(rimg64.FromColor(subim), FGMRConfig(sbin))

	// Take rectangle in f of same size as g.
	min := rect.Min.Div(sbin)
	subf := f.SubImage(image.Rectangle{min, min.Add(g.Size())})

	const prec = 1e-9
	// Skip last element because of a slight difference.
	// (Not using cell outside image.)
	for x := 0; x < g.Width; x++ {
		for y := 0; y < g.Height; y++ {
			for d := 0; d < g.Channels; d++ {
				want := g.At(x, y, d)
				got := subf.At(x, y, d)
				if math.Abs(want-got) > prec {
					t.Errorf("wrong value: at %d, %d, %d: want %g, got %g", x, y, d, want, got)
				}
			}
		}
	}
}

func ensureDivis(src image.Image, m int) image.Image {
	s := src.Bounds().Size()
	s = s.Sub(s.Mod(image.Rect(0, 0, m, m)))
	dst := image.NewRGBA(image.Rectangle{image.ZP, s})
	draw.Draw(dst, dst.Bounds(), src, image.ZP, draw.Src)
	return dst
}

// Output may be one pixel smaller.
func validSize(size, ref image.Point) bool {
	return ((size.X == ref.X || size.X == ref.X-1) &&
		(size.Y == ref.Y || size.Y == ref.Y-1))
}

func BenchmarkGo(b *testing.B) {
	const (
		sbin  = 8
		fname = "000084.jpg"
	)

	// Load image.
	file, err := os.Open(fname)
	if err != nil {
		b.Fatal(err)
	}
	im, _, err := image.Decode(file)
	if err != nil {
		b.Fatal(err)
	}

	f := rimg64.FromColor(im)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HOG(f, FGMRConfig(sbin))
	}
}

func BenchmarkC(b *testing.B) {
	const (
		sbin  = 8
		fname = "000084.jpg"
	)

	// Load image.
	file, err := os.Open(fname)
	if err != nil {
		b.Fatal(err)
	}
	im, _, err := image.Decode(file)
	if err != nil {
		b.Fatal(err)
	}

	f := rimg64.FromColor(im)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fgmr(f, sbin)
	}
}

func TestTransform_MinInputSize(t *testing.T) {
	cases := []struct {
		Bin   int
		Image image.Point
	}{
		{8, image.Pt(140, 160)},
		{8, image.Pt(141, 161)},
		{8, image.Pt(142, 162)},
		{8, image.Pt(143, 163)},
		{8, image.Pt(144, 164)},
		{8, image.Pt(145, 165)},
		{8, image.Pt(146, 166)},
		{8, image.Pt(147, 167)},
		{9, image.Pt(140, 160)},
		{9, image.Pt(141, 161)},
		{9, image.Pt(142, 162)},
		{9, image.Pt(143, 163)},
		{9, image.Pt(144, 164)},
		{9, image.Pt(145, 165)},
		{9, image.Pt(146, 166)},
		{9, image.Pt(147, 167)},
		{9, image.Pt(148, 168)},
	}
	for _, x := range cases {
		phi := Transform{FGMRConfig(x.Bin)}
		origFeat := phi.Size(x.Image)
		minInput := phi.MinInputSize(origFeat)
		minFeat := phi.Size(minInput)
		if !origFeat.Eq(minFeat) {
			t.Errorf("min input size %v for feature size %v: got features of size %v", minInput, origFeat, minFeat)
			continue
		}
		//t.Logf("input %v has same transform %v as %v", minInput, minFeat, origFeat)
		for _, delta := range []image.Point{{1, 0}, {0, 1}, {1, 1}} {
			nextInput := minInput.Sub(delta)
			nextFeat := phi.Size(nextInput)
			if nextFeat.X >= origFeat.X && nextFeat.Y >= origFeat.Y {
				t.Errorf("found smaller input than %v for feature size %v: %v", minInput, origFeat, nextInput)
				break
			}
			//t.Logf("smaller input %v has smaller transform %v than %v", nextInput, nextFeat, origFeat)
		}
	}
}

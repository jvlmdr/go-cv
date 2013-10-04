package hog

import (
	"github.com/jackvalmadre/go-cv/rimg64"

	"image"
	_ "image/jpeg"
	"math"
	"os"
	"testing"
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
	img, _, err := image.Decode(file)
	// Convert to real values.
	f := rimg64.FromColor(img)

	g := HOG(f, sbin)
	ref := FGMR(f, sbin)

	if !g.Size().Eq(ref.Size()) {
		t.Fatalf("different sizes: want %v, got %v", ref.Size(), g.Size())
	}
	if g.Channels != ref.Channels {
		t.Fatalf("different number of channels: want %v, got %v", ref.Channels, g.Channels)
	}

	const eps = 1e-6
	for x := 0; x < ref.Width; x++ {
		for y := 0; y < ref.Height; y++ {
			for d := 0; d < ref.Channels; d++ {
				want := ref.At(x, y, d)
				got := g.At(x, y, d)
				if math.Abs(want-got) > eps {
					t.Errorf("wrong value at (%d, %d, %d): want %g, got %g", x, y, d, want, got)
				}
			}
		}
	}
}

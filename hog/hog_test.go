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

	g := HOG(f, FGMRConfig(sbin))
	ref := FGMR(f, sbin)

	if !validSize(g.Size(), ref.Size()) {
		t.Fatalf("different sizes: want %v, got %v", ref.Size(), g.Size())
	}
	if g.Channels != ref.Channels {
		t.Fatalf("different number of channels: want %v, got %v", ref.Channels, g.Channels)
	}

	const prec = 1e-9
	// Skip last element because of a slight difference.
	// (Not using cell outside image.)
	for x := 0; x < g.Width-1; x++ {
		for y := 0; y < g.Height-1; y++ {
			for d := 0; d < g.Channels; d++ {
				want := ref.At(x, y, d)
				got := g.At(x, y, d)
				if math.Abs(want-got) > prec {
					t.Errorf("wrong value at (%d, %d, %d): want %g, got %g", x, y, d, want, got)
				}
			}
		}
	}
}

// Output may be one pixel smaller.
func validSize(size, ref image.Point) bool {
	return ((size.X == ref.X || size.X == ref.X-1) &&
		(size.Y == ref.Y || size.Y == ref.Y-1))
}

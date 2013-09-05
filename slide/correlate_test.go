package slide

import (
	"github.com/jackvalmadre/go-cv"
	"github.com/jackvalmadre/lin-go/vec"

	"image"
	"math"
	"testing"
)

// Compare naive and Fourier implementations.
func TestCorrelateImages(t *testing.T) {
	const (
		m   = 40
		n   = 30
		w   = 100
		h   = 80
		eps = 1e-12
	)

	f := cv.NewRealImage(w, h)
	g := cv.NewRealImage(m, n)
	// Initialize to random.
	vec.Copy(f.Vec(), vec.Randn(f.Vec().Len()))
	vec.Copy(g.Vec(), vec.Randn(g.Vec().Len()))

	naive := correlateImagesNaive(f, g)
	fourier := correlateImagesFourier(f, g)

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

// Compare naive and Fourier implementations.
func TestCorrelateVectorImages(t *testing.T) {
	const (
		m   = 40
		n   = 30
		w   = 100
		h   = 80
		c   = 8
		eps = 1e-9
	)

	f := cv.NewRealVectorImage(w, h, c)
	g := cv.NewRealVectorImage(m, n, c)
	// Initialize to random.
	vec.Copy(f.Vec(), vec.Randn(f.Vec().Len()))
	vec.Copy(g.Vec(), vec.Randn(g.Vec().Len()))

	naive := correlateVectorImagesNaive(f, g)
	fourier := correlateVectorImagesFourier(f, g)

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

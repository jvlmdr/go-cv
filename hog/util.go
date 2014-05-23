package hog

import (
	"math"

	"github.com/jackvalmadre/go-cv/rimg64"
)

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}

func round(x float64) int {
	if x < 0 {
		return round(-x)
	}
	return int(math.Floor(x + 0.5))
}

func modf(x float64) (int, float64) {
	if x < 0 {
		// Round down not towards zero.
		a, b := modf(-x)
		return -a - 1, -b + 1
	}
	a := int(x)
	b := x - float64(a)
	return a, b
}

// Avoids f.Set(x, y, f.Get(x, y, ...)).
func addTo(f *rimg64.Image, x, y int, v float64) {
	f.Set(x, y, f.At(x, y)+v)
}

// Avoids f.Set(x, y, d, f.Get(x, y, d, ...)).
func addToMulti(f *rimg64.Multi, x, y, d int, v float64) {
	f.Set(x, y, d, f.At(x, y, d)+v)
}

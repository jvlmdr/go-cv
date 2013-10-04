package hog

import (
	"github.com/jackvalmadre/go-cv/rimg64"

	"math"
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
		a, b := modf(-x)
		return -a - 1, 1 - b
	}
	a, b := math.Modf(x)
	return int(a), b
}

// Avoids f.Set(x, y, f.Get(x, y, ...)).
func addTo(f *rimg64.Image, x, y int, v float64) {
	f.Set(x, y, f.At(x, y)+v)
}

// Avoids f.Set(x, y, d, f.Get(x, y, d, ...)).
func addToMulti(f *rimg64.Multi, x, y, d int, v float64) {
	f.Set(x, y, d, f.At(x, y, d)+v)
}

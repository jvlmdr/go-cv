package detect

import (
	"image"
)

func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}

func centroid(r image.Rectangle) (float64, float64) {
	return float64(r.Min.X+r.Max.X) / 2, float64(r.Min.Y+r.Max.Y) / 2
}

func round(x float64) int {
	if x < 0 {
		return -int(-x + 0.5) // -round(-x)
	}
	return int(x + 0.5)
}

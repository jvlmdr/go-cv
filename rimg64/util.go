package rimg64

import "math"

func round(x float64) int {
	return int(math.Floor(x + 0.5))
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

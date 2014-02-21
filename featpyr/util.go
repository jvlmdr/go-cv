package featpyr

import "math"

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// Returns the nearest integer.
func round(x float64) int {
	return int(math.Floor(x + 0.5))
}

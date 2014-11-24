package slide

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

func logb(x int) int {
	n := 0
	for x != 0 {
		x /= 2
		n++
	}
	return n
}

func sqr(x float64) float64 { return x * x }

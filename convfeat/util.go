package convfeat

func sqr(x float64) float64 {
	return x * x
}

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

func ceilDiv(a, b int) int {
	return (a + b - 1) / b
}

package slide

import "image"

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

func ceilDiv(a, b int) int {
	return (a + b - 1) / b
}

func ceilDivPt(a image.Point, b int) image.Point {
	return image.Pt(ceilDiv(a.X, b), ceilDiv(a.Y, b))
}

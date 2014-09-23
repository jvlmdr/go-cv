package slide

import "image"

// ValidSize returns the number of positions such that
// the template g lies entirely inside the image f.
func ValidSize(f, g image.Point) image.Point {
	var h image.Point
	h.X = max(f.X-g.X+1, 0)
	h.Y = max(f.Y-g.Y+1, 0)
	return h
}

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

// ValidRect returns the region of non-periodic convolution
// (or correlation) of a template g with an image f.
// It is assumed that both are packed in the top left corner.
func ValidRect(f, g image.Point, corr bool) image.Rectangle {
	// Compute size of region.
	var s image.Point
	s.X = max(f.X-g.X+1, 0)
	s.Y = max(f.Y-g.Y+1, 0)
	r := image.Rectangle{Max: s}
	if corr {
		return r
	}
	return r.Add(image.Pt(g.X-1, g.Y-1))
}

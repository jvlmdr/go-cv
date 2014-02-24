package detect

import "image"

func iou(p, q image.Rectangle) float64 {
	inter := area(p.Intersect(q))
	union := area(p) + area(q) - inter
	return float64(inter) / float64(union)
}

func interRel(a, b image.Rectangle) float64 {
	inter := area(a.Intersect(b))
	return float64(inter) / float64(area(a))
}

func area(r image.Rectangle) int {
	s := r.Size()
	return s.X * s.Y
}

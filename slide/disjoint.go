package slide

import (
	"container/list"
	"image"
	"sort"

	"github.com/jackvalmadre/go-cv/rimg64"
)

// Find detections which do not intersect
// given the response to and size of a template.
func Disjoint(f *rimg64.Image, size image.Point) []image.Point {
	var pixels []image.Point
	for x := 0; x < f.Width; x++ {
		for y := 0; y < f.Height; y++ {
			pixels = append(pixels, image.Pt(x, y))
		}
	}
	sort.Sort(sort.Reverse(pixelsByScore{f, pixels}))

	// Remaining pixels in image.
	rem := list.New()
	// Point to position of elements in list.
	elem := make([][]*list.Element, f.Width)
	for i := range elem {
		elem[i] = make([]*list.Element, f.Height)
	}
	// Populate structures.
	for _, p := range pixels {
		elem[p.X][p.Y] = rem.PushBack(p)
	}

	var detections []image.Point
	for rem.Len() > 0 {
		p, ok := rem.Front().Value.(image.Point)
		if !ok {
			panic("Invalid type")
		}
		detections = append(detections, p)

		a := max(0, p.X-size.X+1)
		b := min(f.Width-1, p.X+size.X-1)
		for x := a; x <= b; x++ {
			c := max(0, p.Y-size.Y+1)
			d := min(f.Height-1, p.Y+size.Y-1)
			for y := c; y <= d; y++ {
				if elem[x][y] == nil {
					continue
				}
				rem.Remove(elem[x][y])
				elem[x][y] = nil
			}
		}
	}

	return detections
}

type pixelsByScore struct {
	Image  *rimg64.Image
	Pixels []image.Point
}

func (s pixelsByScore) Len() int { return len(s.Pixels) }

func (s pixelsByScore) Less(i, j int) bool {
	p := s.Pixels[i]
	q := s.Pixels[j]
	return s.Image.At(p.X, p.Y) < s.Image.At(q.X, q.Y)
}

func (s pixelsByScore) Swap(i, j int) {
	s.Pixels[i], s.Pixels[j] = s.Pixels[j], s.Pixels[i]
}

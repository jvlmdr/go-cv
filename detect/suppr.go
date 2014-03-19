package detect

import (
	"container/list"
	"image"
	"sort"
)

// Determines whether a overlaps b.
// The score of a was higher than that of b.
type OverlapFunc func(a, b image.Rectangle) bool

// Intersection over union.
func IOU(a, b image.Rectangle) float64 {
	inter := area(a.Intersect(b))
	union := area(a) + area(b) - inter
	return float64(inter) / float64(union)
}

// Returns the fraction of B which is covered by A.
func Cover(a, b image.Rectangle) float64 {
	inter := area(a.Intersect(b))
	return float64(inter) / float64(area(b))
}

func area(r image.Rectangle) int {
	s := r.Size()
	return s.X * s.Y
}

// Returns a list of indices to keep.
//
//	// Two rectangles overlap if their IOU exceeds 0.5.
//	overlap := func(a, b image.Rectangle) bool { return detect.IOU(a, b) >= 0.5 }
//	Suppress(dets, maxnum, overlap)
func Suppress(dets []Det, maxnum int, overlap OverlapFunc) []int {
	if !sort.IsSorted(detsByScoreDesc(dets)) {
		panic("not sorted")
	}
	// Copy into linked list.
	rem := list.New()
	for i := range dets {
		rem.PushBack(i)
	}
	// Select best detection, remove those which overlap with it.
	var subset []int
	for rem.Len() > 0 && (maxnum <= 0 || len(subset) < maxnum) {
		subset = append(subset, pop(rem, dets, overlap))
	}
	return subset
}

func SuppressDets(dets []Det, maxnum int, overlap OverlapFunc) []Det {
	inds := Suppress(dets, maxnum, overlap)
	subset := make([]Det, len(inds))
	for i, ind := range inds {
		subset[i] = dets[ind]
	}
	return subset
}

func pop(rem *list.List, dets []Det, overlap OverlapFunc) int {
	i := rem.Remove(rem.Front()).(int)
	var next *list.Element
	for e := rem.Front(); e != nil; e = next {
		// Buffer next so that we can remove e.
		next = e.Next()
		// Get candidate detection.
		j := e.Value.(int)
		// Suppress if the rectangles are deemed to overlap.
		// The first argument has the higher score.
		if overlap(dets[i].Rect, dets[j].Rect) {
			// Remove.
			rem.Remove(e)
		}
	}
	return i
}

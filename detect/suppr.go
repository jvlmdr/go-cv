package detect

import (
	"container/list"
	"image"
	"sort"
)

func SuppressOverlap(cands []Det, maxnum int, maxinter float64) []Det {
	if !sort.IsSorted(sort.Reverse(byScore(cands))) {
		panic("detections are not sorted descending by score")
	}
	// Copy into linked list.
	rem := list.New()
	for _, det := range cands {
		rem.PushBack(det)
	}
	// Select best detection, remove those which overlap with it.
	var dets []Det
	for rem.Len() > 0 && (maxnum <= 0 || len(dets) < maxnum) {
		det := pop(rem, maxinter)
		dets = append(dets, det)
	}
	return dets
}

func pop(rem *list.List, maxinter float64) Det {
	det := rem.Remove(rem.Front()).(Det)
	var next *list.Element
	for e := rem.Front(); e != nil; e = next {
		// Buffer next so that we can remove e.
		next = e.Next()
		// Get candidate detection.
		cand := e.Value.(Det)
		if intersect(det.Pos, cand.Pos, maxinter) {
			// Remove.
			rem.Remove(e)
		}
	}
	return det
}

func intersect(a, b image.Rectangle, maxinter float64) bool {
	ab := a.Intersect(b)
	rela := float64(area(ab)) / float64(area(a))
	relb := float64(area(ab)) / float64(area(b))
	return (rela > maxinter && relb > maxinter)
}

func area(r image.Rectangle) int {
	s := r.Size()
	return s.X * s.Y
}

type byScore []Det

func (s byScore) Len() int           { return len(s) }
func (s byScore) Less(i, j int) bool { return s[i].Score < s[j].Score }
func (s byScore) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

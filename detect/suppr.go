package detect

import (
	"container/list"
	"image"
	"sort"
)

// Returns a list of indices to keep.
func Suppress(dets []Det, maxnum int, maxinter float64) []int {
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
		subset = append(subset, pop(rem, dets, maxinter))
	}
	return subset
}

func SuppressDets(dets []Det, maxnum int, maxinter float64) []Det {
	inds := Suppress(dets, maxnum, maxinter)
	subset := make([]Det, len(inds))
	for i, ind := range inds {
		subset[i] = dets[ind]
	}
	return subset
}

func pop(rem *list.List, dets []Det, maxinter float64) int {
	i := rem.Remove(rem.Front()).(int)
	var next *list.Element
	for e := rem.Front(); e != nil; e = next {
		// Buffer next so that we can remove e.
		next = e.Next()
		// Get candidate detection.
		j := e.Value.(int)
		// Suppress if both windows overlap each other.
		if interRelBoth(dets[i].Rect, dets[j].Rect, maxinter) {
			// Remove.
			rem.Remove(e)
		}
	}
	return i
}

func interRelBoth(a, b image.Rectangle, maxinter float64) bool {
	return interRel(a, b) > maxinter && interRel(b, a) > maxinter
}

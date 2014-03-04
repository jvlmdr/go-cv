package detect

import (
	"container/list"
	"image"
	"sort"
)

// Validated detection.
type ValDet struct {
	// Need score to merge results from multiple images.
	Det
	// If this was a true detection, give its reference.
	True bool
}

func MergeValDets(a, b []ValDet) []ValDet {
	if !sort.IsSorted(valDetsByScoreDesc(a)) {
		panic("not sorted")
	}
	if !sort.IsSorted(valDetsByScoreDesc(b)) {
		panic("not sorted")
	}

	c := make([]ValDet, 0, len(a)+len(b))
	for len(a) > 0 || len(b) > 0 {
		if len(a) == 0 {
			return append(c, b...)
		}
		if len(b) == 0 {
			return append(c, a...)
		}
		if a[0].Score > b[0].Score {
			c, a = append(c, a[0]), a[1:]
		} else {
			c, b = append(c, b[0]), b[1:]
		}
	}
	return c
}

// Takes a list of detections ordered by score
// and an unordered list of ground truth regions.
// Assigns each detection to the reference window with which it overlaps the most
// if that overlap is sufficient.
//
// Overlap is evaluated using intersection over union.
func ValidateMatch(dets []Det, refs []image.Rectangle, mininter float64) *ResultSet {
	m := Match(dets, refs, mininter)
	return ResultsMatch(dets, refs, m)
}

// Takes a list of detections ordered by score
// and an unordered list of ground truth regions.
// Assigns each detection to the reference window with which it overlaps the most
// if that overlap is sufficient.
// Returns a map from detection index to reference index.
//
// Overlap is evaluated using intersection over union.
func Match(dets []Det, refs []image.Rectangle, mininter float64) map[int]int {
	// Map from dets to refs.
	m := make(map[int]int)
	// List of indices remaining in refs.
	r := list.New()
	for j := range refs {
		r.PushBack(j)
	}

	for i, det := range dets {
		var maxinter float64
		var argmax *list.Element
		for e := r.Front(); e != nil; e = e.Next() {
			j := e.Value.(int)
			inter := iou(det.Rect, refs[j])
			if inter < mininter {
				continue
			}
			if inter > maxinter {
				maxinter = inter
				argmax = e
			}
		}
		// Sufficient overlap with none.
		if argmax == nil {
			continue
		}
		// Found a match.
		m[i] = r.Remove(argmax).(int)
	}
	return m
}

func ResultsMatch(dets []Det, refs []image.Rectangle, m map[int]int) *ResultSet {
	// Label each detection as true positive or false positive.
	valdets := make([]ValDet, len(dets))
	// Record which references were matched.
	used := make(map[int]bool)
	for i, det := range dets {
		j, p := m[i]
		if !p {
			valdets[i] = ValDet{det, false}
			continue
		}
		if used[j] {
			panic("already matched")
		}
		used[j] = true
		valdets[i] = ValDet{det, true}
	}
	misses := len(refs) - len(used)
	return &ResultSet{valdets, misses}
}

type valDetsByScoreDesc []ValDet

func (s valDetsByScoreDesc) Len() int           { return len(s) }
func (s valDetsByScoreDesc) Less(i, j int) bool { return s[i].Score > s[j].Score }
func (s valDetsByScoreDesc) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

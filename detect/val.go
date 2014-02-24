package detect

import (
	"container/list"
	"image"
)

// Validated detection.
type ValDet struct {
	// Need score to merge results from multiple images.
	Det
	// If this was a true detection, give its reference.
	True bool
}

// Can be used to generate an ROC curve.
// Multiple can be merged.
type ResultSet struct {
	// Validated detections ordered by score.
	Dets []ValDet
	// Number of missed instances (false negatives).
	Misses int
}

func (s *ResultSet) Merge(t *ResultSet) *ResultSet {
	dets := mergeValDets(s.Dets, t.Dets)
	misses := s.Misses + t.Misses
	return &ResultSet{dets, misses}
}

func mergeValDets(a, b []ValDet) []ValDet {
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
// Returns a map from detection index to reference index.
//
// Overlap is evaluated using intersection over union.
func Validate(dets []Det, refs []image.Rectangle, mininter float64) map[int]int {
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

func Results(dets []Det, refs []image.Rectangle, m map[int]int) ResultSet {
	// Label each detection as true positive or false positive.
	valdets := make([]ValDet, len(dets))
	// Record which references were matched.
	var used map[int]bool
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
	return ResultSet{valdets, misses}
}

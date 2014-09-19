package detect

import (
	"sort"

	"github.com/jvlmdr/go-ml/ml"
)

// ValScore describes the score and correctness of a detection.
type ValScore struct {
	Score float64
	// Is it a true positive?
	True bool
}

// ValSet summarizes the validation of detections in a set of images.
type ValSet struct {
	// Validated detections ordered (descending) by score.
	Dets []ValScore
	// Number of missed instances (false negatives).
	Misses int
	// Number of images over which this was computed.
	Images int
}

// Clone creates a copy.
func (s *ValSet) Clone() *ValSet {
	if s == nil {
		return nil
	}
	dets := make([]ValScore, len(s.Dets))
	copy(dets, s.Dets)
	return &ValSet{dets, s.Misses, s.Images}
}

// Merge combines the statistics of two sets.
func (s *ValSet) Merge(t *ValSet) *ValSet {
	if t == nil {
		return s.Clone()
	}
	if s == nil {
		return t.Clone()
	}
	dets := mergeScores(s.Dets, t.Dets)
	return &ValSet{dets, s.Misses + t.Misses, s.Images + t.Images}
}

// MergeValSets combines multiple validated sets.
// Use set.Merge() to merge two sets.
func MergeValSets(sets ...*ValSet) *ValSet {
	var (
		dets   []ValScore
		misses int
		images int
	)
	for _, set := range sets {
		dets = append(dets, set.Dets...)
		misses += set.Misses
		images += set.Images
	}
	sort.Sort(valScoresDesc(dets))
	return &ValSet{dets, misses, images}
}

// Enum enumerates the performance of operating points achieved
// by varying the threshold.
func (s *ValSet) Enum() ml.PerfPath {
	if !sort.IsSorted(valScoresDesc(s.Dets)) {
		panic("not sorted")
	}
	results := make([]ml.Perf, len(s.Dets)+1)
	// Start with everything classified negative.
	// TP = 0, FP = 0.
	// FN is number of actual positives.
	// TN not computed, depends on number of windows.
	numPos := numTrue(s.Dets) + s.Misses
	curr := ml.Perf{TP: 0, FP: 0, FN: numPos}
	results[0] = curr
	// Set each detection to positive in order.
	for i, det := range s.Dets {
		if det.True {
			curr.TP++
			curr.FN--
		} else {
			curr.FP++
		}
		results[i+1] = curr
	}
	return results
}

func numTrue(scores []ValScore) int {
	var n int
	for _, s := range scores {
		if s.True {
			n++
		}
	}
	return n
}

func mergeScores(a, b []ValScore) []ValScore {
	if !sort.IsSorted(valScoresDesc(a)) {
		panic("not sorted")
	}
	if !sort.IsSorted(valScoresDesc(b)) {
		panic("not sorted")
	}

	c := make([]ValScore, 0, len(a)+len(b))
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

type valScoresDesc []ValScore

func (s valScoresDesc) Len() int           { return len(s) }
func (s valScoresDesc) Less(i, j int) bool { return s[i].Score > s[j].Score }
func (s valScoresDesc) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

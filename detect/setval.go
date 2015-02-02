package detect

import (
	"fmt"
	"log"
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

// MissRateAtFPPIs computes the miss rates at multiple
// false-positive-per-image rates.
// An error results if the number of false positives to
// achieve some FPPI is not greater than zero.
func MissRateAtFPPIs(valset *ValSet, fppis []float64) ([]float64, error) {
	// Construct cumulative count of true positives
	// as threshold decreases.
	n := len(valset.Dets)
	falsePos := make([]int, n+1)
	for i, det := range valset.Dets {
		if det.True {
			falsePos[i+1] = falsePos[i]
		} else {
			falsePos[i+1] = falsePos[i] + 1
		}
	}
	// Number of detections that were positive at any threshold.
	var posDets int
	for _, det := range valset.Dets {
		if det.True {
			posDets++
		}
	}

	missRates := make([]float64, len(fppis))
	for i, fppi := range fppis {
		// Find operating point at which false positives
		// are not more than number given.
		// Obtain absolute number of false positives.
		// Largest integer such that maxFalsePos / numImages <= fppi.
		maxFalsePos := int(fppi * float64(valset.Images))
		// This may be zero, in which case there are
		// not enough images to measure miss rate at this FPPI.
		if maxFalsePos < 1 {
			return nil, fmt.Errorf("not enough images: fppi %g, images %d", fppi, valset.Images)
		}
		// Note that if the false positive is followed by true positives,
		// then there will be several points with the same FPPI.
		// falsePos[i] is the number of false positives in Dets[0, ..., i-1].
		// Want to find the largest i such that
		//   falsePos[i] <= maxFalsePos
		// and then use detections from Dets[0, ..., i-1].
		// This is equivalent to the smallest i such that
		//   falsePos[i+1] > maxFalsePos.
		// Check upper boundary:
		// The greatest i that Search() will test is i = n-1:
		//   falsePos[n] > maxFalsePos.
		// We should use Dets[0, ..., n-1] if
		//   falsePos[n] <= maxFalsePos.
		l := sort.Search(n, func(i int) bool { return falsePos[i+1] > maxFalsePos })

		// Number of true detections in Dets[0, ..., l-1].
		truePos := l - falsePos[l]
		// Number of true detections in Dets[l, ..., n-1].
		falseNeg := posDets - truePos
		// Miss rate is false negatives / number of actual positives.
		missRate := float64(falseNeg+valset.Misses) / float64(posDets+valset.Misses)
		log.Printf("fppi %.3g, false pos %d, num dets %d, miss rate %.3g", fppi, maxFalsePos, l, missRate)
		missRates[i] = missRate
	}
	return missRates, nil
}

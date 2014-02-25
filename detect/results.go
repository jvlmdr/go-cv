package detect

import "github.com/jackvalmadre/go-ml"

// Can be used to generate an ROC curve.
// Multiple can be merged.
type ResultSet struct {
	// Validated detections ordered by score.
	Dets []ValDet
	// Number of missed instances (false negatives).
	Misses int
}

func (s *ResultSet) Merge(t *ResultSet) *ResultSet {
	if t == nil {
		return s.Clone()
	}
	if s == nil {
		return t.Clone()
	}
	dets := MergeValDets(s.Dets, t.Dets)
	misses := s.Misses + t.Misses
	return &ResultSet{dets, misses}
}

func (s *ResultSet) Clone() *ResultSet {
	if s == nil {
		return nil
	}
	dets := make([]ValDet, len(s.Dets))
	copy(dets, s.Dets)
	return &ResultSet{dets, s.Misses}
}

func (s *ResultSet) Enum() []ml.Result {
	n := len(s.Dets)
	results := make([]ml.Result, n+1)
	// Start with everything classified negative.
	// TP = 0, FP = 0.
	// FN (number of actual positives) is unknown.
	// TN not computed, depends on number of windows.
	var curr ml.Result
	results[0] = curr

	// Set each detection to positive in order.
	for i, det := range s.Dets {
		if det.True {
			curr.TP++
		} else {
			curr.FP++
		}
		results[i+1] = curr
	}

	// Now the number of actual positives is known.
	pos := results[n].TP + s.Misses
	// Ensure that TP + FN = constant.
	for i := range results {
		results[i].FN = pos - results[i].TP
	}
	return results
}
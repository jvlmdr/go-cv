package detect

import (
	"container/list"
	"image"
)

type Val struct {
	// Is this a true detection?
	True bool
	// If so, give its corresponding annotation.
	Ref image.Rectangle
}

// ValDet describes a validated detection,
// a detection that has been compared to ground truth and
// deemed true or false.
type ValDet struct {
	Det
	Val
}

// ValImage describes the validation of a set of detections in one image.
// Can be used to visualize correct detections, false positives and missed detections.
type ValImage struct {
	// Validated detections ordered (descending) by score.
	Dets []ValDet
	// Number of missed instances (false negatives).
	Misses []image.Rectangle
}

// Scores extracts just the scores of the validated detections,
// discarding their location.
func (im *ValImage) Scores() []ValScore {
	scores := make([]ValScore, len(im.Dets))
	for i, det := range im.Dets {
		scores[i] = ValScore{det.Score, det.True}
	}
	return scores
}

// Set returns the set containing this image.
func (im *ValImage) Set() *ValSet {
	return &ValSet{im.Scores(), len(im.Misses), 1}
}

// Validate compares a list of detections in an image to a ground-truth reference.
// The detections must be ordered (descending) by score.
// Some of the regions can be labelled as "ignore".
// Each detection is assigned to the reference window with which it overlaps the most
// if that overlap is sufficient.
// This is performed greedily and references cannot match to multiple detections.
//
// Sufficient overlap to match a reference is assessed using intersection-over-union.
// Sufficient overlap to ignore a detection is assessed by what fraction of the detection is covered.
func Validate(dets []Det, refs, ignore []image.Rectangle, refMinIOU, ignoreMinCover float64) *ValImage {
	vals, miss := ValidateList(DetSlice(dets), refs, ignore, refMinIOU, ignoreMinCover)
	valdets := make([]ValDet, len(dets))
	for i := range dets {
		valdets[i] = ValDet{dets[i], vals[i]}
	}
	return &ValImage{valdets, miss}
}

func ValidateList(dets DetList, refs, ignore []image.Rectangle, refMinIOU, ignoreMinCover float64) (vals []Val, miss []image.Rectangle) {
	// Match rectangles and then remove any detections (incorrect or otherwise) which are ignored.
	m := Match(dets, refs, refMinIOU)
	// Label each detection as true positive or false positive.
	vals = make([]Val, dets.Len())
	// Record which references were matched.
	used := make(map[int]bool)
	for i := 0; i < dets.Len(); i++ {
		det := dets.At(i)
		j, p := m[i]
		if !p {
			// Detection did not have a match.
			// Check whether to ignore the false positive.
			if anyCovers(ignore, det.Rect, ignoreMinCover) {
				continue
			}
			vals[i] = Val{True: false}
			continue
		}
		if used[j] {
			panic("already matched")
		}
		used[j] = true
		vals[i] = Val{True: true, Ref: refs[j]}
	}
	miss = make([]image.Rectangle, 0, len(refs)-len(used))
	for j, ref := range refs {
		if used[j] {
			continue
		}
		miss = append(miss, ref)
	}
	return vals, miss
}

// Match takes a list of detections ordered (descending) by score
// and an unordered list of ground truth regions.
// Assigns each detection to the reference window with which it overlaps the most
// if that overlap is sufficient.
// Returns a map from detection index to reference index.
//
// Overlap is evaluated using intersection over union.
// The scores of the detections are not used.
func Match(dets DetList, refs []image.Rectangle, mininter float64) map[int]int {
	// Map from dets to refs.
	m := make(map[int]int)
	// List of indices remaining in refs.
	r := list.New()
	for j := range refs {
		r.PushBack(j)
	}

	for i := 0; i < dets.Len(); i++ {
		det := dets.At(i)
		var maxinter float64
		var argmax *list.Element
		for e := r.Front(); e != nil; e = e.Next() {
			j := e.Value.(int)
			inter := IOU(det.Rect, refs[j])
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

func anyCovers(ys []image.Rectangle, x image.Rectangle, minCover float64) bool {
	for _, y := range ys {
		if covers(y, x, minCover) {
			return true
		}
	}
	return false
}

// Computes whether A covers B.
func covers(a, b image.Rectangle, min float64) bool {
	return float64(area(b.Intersect(a)))/float64(area(b)) > min
}

type valDetsByScoreDesc []ValDet

func (s valDetsByScoreDesc) Len() int           { return len(s) }
func (s valDetsByScoreDesc) Less(i, j int) bool { return s[i].Score > s[j].Score }
func (s valDetsByScoreDesc) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

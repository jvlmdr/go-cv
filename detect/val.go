package detect

import (
	"container/list"
	"image"
)

// ValDet describes a validated detection,
// a detection that has been compared to ground truth and
// deemed true or false.
type ValDet struct {
	Det
	// Is this a true detection?
	True bool
	// If so, give its corresponding annotation.
	Ref image.Rectangle
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
	im := matchesToValImage(Match(dets, refs, refMinIOU), dets, refs)
	if len(ignore) > 0 {
		im.Dets = removeIfIgnored(im.Dets, ignore, ignoreMinCover)
	}
	return im
}

// Match takes a list of detections ordered (descending) by score
// and an unordered list of ground truth regions.
// Assigns each detection to the reference window with which it overlaps the most
// if that overlap is sufficient.
// Returns a map from detection index to reference index.
//
// Overlap is evaluated using intersection over union.
// The scores of the detections are not used.
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

func matchesToValImage(m map[int]int, dets []Det, refs []image.Rectangle) *ValImage {
	// Label each detection as true positive or false positive.
	valdets := make([]ValDet, len(dets))
	// Record which references were matched.
	used := make(map[int]bool)
	for i, det := range dets {
		j, p := m[i]
		if !p {
			valdets[i] = ValDet{Det: det, True: false}
			continue
		}
		if used[j] {
			panic("already matched")
		}
		used[j] = true
		valdets[i] = ValDet{Det: det, True: true, Ref: refs[j]}
	}
	misses := make([]image.Rectangle, 0, len(refs)-len(used))
	for j, ref := range refs {
		if used[j] {
			continue
		}
		misses = append(misses, ref)
	}
	return &ValImage{valdets, misses}
}

func removeIfIgnored(dets []ValDet, ignore []image.Rectangle, minCover float64) []ValDet {
	var keep []ValDet
	for _, det := range dets {
		if !det.True && anyCovers(ignore, det.Rect, minCover) {
			continue
		}
		keep = append(keep, det)
	}
	return keep
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

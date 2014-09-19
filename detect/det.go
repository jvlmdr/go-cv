package detect

import (
	"image"
	"math"
	"sort"
)

// Det describes one detection.
type Det struct {
	Score float64
	Rect  image.Rectangle
}

// Sort sorts a list of detections descending by score.
func Sort(dets []Det) {
	if anyIsNaN(dets) {
		panic("cannot sort scores: NaN")
	}
	sort.Sort(detsByScoreDesc(dets))
}

type detsByScoreDesc []Det

func (s detsByScoreDesc) Len() int      { return len(s) }
func (s detsByScoreDesc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s detsByScoreDesc) Less(i, j int) bool {
	return s[i].Score > s[j].Score
}

func anyIsNaN(dets []Det) bool {
	for _, det := range dets {
		if math.IsNaN(det.Score) {
			return true
		}
	}
	return false
}

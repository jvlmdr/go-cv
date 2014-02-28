package detect

import (
	"image"
	"sort"
)

// Detection in image.
type Det struct {
	Score float64
	Rect  image.Rectangle
}

func MergeDets(a, b []Det) []Det {
	c := make([]Det, 0, len(a)+len(b))
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

func SortDets(dets []Det) {
	sort.Sort(byScoreDesc(dets))
}

type byScoreDesc []Det

func (s byScoreDesc) Len() int           { return len(s) }
func (s byScoreDesc) Less(i, j int) bool { return s[i].Score > s[j].Score }
func (s byScoreDesc) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

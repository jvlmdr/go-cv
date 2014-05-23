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

func MergeTwoDets(a, b []Det) []Det {
	if !sort.IsSorted(detsByScoreDesc(a)) {
		panic("not sorted")
	}
	if !sort.IsSorted(detsByScoreDesc(b)) {
		panic("not sorted")
	}

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

func cloneDets(src []Det) []Det {
	dst := make([]Det, len(src))
	copy(dst, src)
	return dst
}

func MergeDets(dets ...[]Det) []Det {
	switch len(dets) {
	case 0:
		return nil
	case 1:
		return cloneDets(dets[0])
	case 2:
		return MergeTwoDets(dets[0], dets[1])
	}

	var all []Det
	for _, d := range dets {
		all = append(all, d...)
	}
	sort.Sort(detsByScoreDesc(all))
	return all
}

func Sort(dets []Det) {
	sort.Sort(detsByScoreDesc(dets))
}

type detsByScoreDesc []Det

func (s detsByScoreDesc) Len() int           { return len(s) }
func (s detsByScoreDesc) Less(i, j int) bool { return s[i].Score > s[j].Score }
func (s detsByScoreDesc) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

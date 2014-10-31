package exemplar

import (
	"image"

	"github.com/jvlmdr/go-cv/detect"
)

type ValDet struct {
	Det
	detect.Val
}

type ValImage struct {
	Dets   []ValDet
	Misses []image.Rectangle
}

func Validate(dets []Det, refs, ignore []image.Rectangle, refMinIOU, ignoreMinCover float64) *ValImage {
	vals, miss := detect.ValidateList(DetSlice(dets), refs, ignore, refMinIOU, ignoreMinCover)
	valdets := make([]ValDet, len(dets))
	for i := range dets {
		valdets[i] = ValDet{dets[i], vals[i]}
	}
	return &ValImage{valdets, miss}
}

func (im *ValImage) Scores() []detect.ValScore {
	scores := make([]detect.ValScore, len(im.Dets))
	for i, det := range im.Dets {
		scores[i] = detect.ValScore{det.Score, det.True}
	}
	return scores
}

func (im *ValImage) Set() *detect.ValSet {
	return &detect.ValSet{im.Scores(), len(im.Misses), 1}
}

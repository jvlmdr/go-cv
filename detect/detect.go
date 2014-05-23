package detect

import (
	"image"
	"sort"

	"github.com/jackvalmadre/go-cv/feat"
	"github.com/jackvalmadre/go-cv/rimg64"
	"github.com/jackvalmadre/go-cv/slide"
)

// Performs detection and non-max suppression.
// Returns a list of scored detection windows.
// Windows are specified as rectangles in the original pixel image.
//
// It is necessary to provide:
// the integer downsample rate of the feature transform,
// the margin which was added to the image before taking the feature transform,
// the rectangular region within the window which corresponds to the annotation.
func Detect(im *rimg64.Multi, margin feat.Margin, tmpl *FeatTmpl, rate int, detopts DetFilter, suppropts SupprFilter) []Det {
	// Evaluate detector and extract rectangles.
	pts := Points(im, tmpl.Image, detopts.LocalMax, detopts.MinScore)
	dets := make([]Det, len(pts))
	for i, det := range pts {
		rect := FeatPtToImRect(det.Point, rate, margin, tmpl.Interior)
		dets[i] = Det{det.Score, rect}
	}
	// Perform non-max suppression.
	Sort(dets)
	dets = Suppress(dets, suppropts.MaxNum, suppropts.Overlap)
	return dets
}

type DetFilter struct {
	// Ignore detections which are smaller than a neighbor?
	LocalMax bool
	// Score threshold.
	MinScore float64
}

// Describes a scored position.
type DetPos struct {
	Score float64
	image.Point
}

// Returns a list of scored detection windows.
// Windows are described by their top-left corner in the feature image.
//
// If localmax is true, then points which have a neighbor greater than them are excluded.
// Any windows less than minscore are excluded.
func Points(im *rimg64.Multi, tmpl *rimg64.Multi, localmax bool, minscore float64) []DetPos {
	resp := slide.CorrMulti(im, tmpl)
	var dets []DetPos
	// Iterate over positions and check criteria.
	for u := 0; u < resp.Width; u++ {
		for v := 0; v < resp.Height; v++ {
			score := resp.At(u, v)
			if score < minscore {
				continue
			}
			if localmax && !localMax(resp, u, v) {
				continue
			}
			// Scale by rate, then apply margin and interior offsets.
			dets = append(dets, DetPos{score, image.Pt(u, v)})
		}
	}
	return dets
}

// Converts the position of a detection in a feature image to a rectangle in the intensity image.
//
// Additional arguments are:
// the integer downsample rate of the feature transform,
// the margin which was added to the image before taking the feature transform,
// the rectangular region within the window which corresponds to the annotation.
func FeatPtToImRect(pt image.Point, rate int, margin feat.Margin, interior image.Rectangle) image.Rectangle {
	return interior.Add(pt.Mul(rate)).Sub(margin.TopLeft())
}

func localMax(r *rimg64.Image, u, v int) bool {
	uv := r.At(u, v)
	if u > 0 && r.At(u-1, v) > uv {
		return false
	}
	if u < r.Width-1 && r.At(u+1, v) > uv {
		return false
	}
	if v > 0 && r.At(u, v-1) > uv {
		return false
	}
	if v < r.Height-1 && r.At(u, v+1) > uv {
		return false
	}
	return true
}

func sortPoints(resp *rimg64.Image, pts []image.Point) {
	sort.Sort(ptsByScore{resp, pts})
}

type ptsByScore struct {
	Resp  *rimg64.Image
	Elems []image.Point
}

func (s ptsByScore) Len() int { return len(s.Elems) }

func (s ptsByScore) Less(i, j int) bool {
	p := s.Elems[i]
	q := s.Elems[j]
	return s.Resp.At(p.X, p.Y) > s.Resp.At(q.X, q.Y)
}

func (s ptsByScore) Swap(i, j int) { s.Elems[i], s.Elems[j] = s.Elems[j], s.Elems[i] }

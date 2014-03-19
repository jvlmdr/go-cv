package featpyr

import (
	"github.com/jackvalmadre/go-cv/detect"
	"github.com/jackvalmadre/go-cv/imgpyr"
	"github.com/jackvalmadre/go-cv/rimg64"
	"github.com/jackvalmadre/go-cv/slide"

	"image"
	"log"
	"sort"
)

// Detection in pyramid.
type Det struct {
	Score float64
	// Location in feature pyramid. Not pixel pyramid!
	imgpyr.Point
}

// Searches each level of the pyramid with a sliding window.
// Returns the detections ordered by score, from highest to lowest.
//
// The size of the template in pixels must be provided.
func DetectPoints(pyr *Pyramid, tmpl *detect.FeatTmpl, overlap detect.OverlapFunc, minscore float64, maxnum int, localmax bool) []Det {
	resps := evalTmpl(pyr, tmpl.Image)
	pts := filterDets(resps, minscore, localmax)
	// Sort then convert to rectangles.
	sort.Sort(byScoreDesc(pts))
	rects := pointsToRects(pyr, pts, tmpl.Interior)
	// Discard rectangles and take order.
	order := detect.Suppress(rects, maxnum, overlap)
	for i := range order {
		pts[i] = pts[order[i]]
	}
	pts = pts[:len(order)]
	return pts
}

// Searches each level of the pyramid with a sliding window.
// Returns the detections ordered by score, from highest to lowest.
//
// The size of the template in pixels must be provided.
func Detect(pyr *Pyramid, tmpl *detect.FeatTmpl, overlap detect.OverlapFunc, minscore float64, maxnum int, localmax bool) []detect.Det {
	resps := evalTmpl(pyr, tmpl.Image)
	pts := filterDets(resps, minscore, localmax)
	// Convert to rectangles then sort.
	rects := pointsToRects(pyr, pts, tmpl.Interior)
	detect.SortDets(rects)
	return detect.SuppressDets(rects, maxnum, overlap)
}

func pointsToRects(pyr *Pyramid, pts []Det, interior image.Rectangle) []detect.Det {
	dets := make([]detect.Det, len(pts))
	for i, pt := range pts {
		rect := pyr.ToImageRect(pt.Point, interior)
		dets[i] = detect.Det{pt.Score, rect}
	}
	return dets
}

func evalTmpl(pyr *Pyramid, tmpl *rimg64.Multi) []*rimg64.Image {
	// Evaluate sliding window at each level.
	resps := make([]*rimg64.Image, 0, len(pyr.Feats))
	for _, feat := range pyr.Feats {
		width := feat.Width - tmpl.Width + 1
		height := feat.Height - tmpl.Height + 1
		if width < 1 || height < 1 {
			break
		}
		resps = append(resps, slide.CorrMulti(feat, tmpl))
	}
	log.Print("evaluated template across all levels")
	return resps
}

func filterDets(resps []*rimg64.Image, minscore float64, localmax bool) []Det {
	// Evaluate sliding window at each level.
	var dets []Det
	if localmax {
		dets = peaks(resps, minscore)
	} else {
		// Put every window into a list.
		dets = windows(resps, minscore)
	}
	return dets
}

// Extracts position and score per pixel in the response images.
func windows(resps []*rimg64.Image, minscore float64) []Det {
	var dets []Det
	for k, resp := range resps {
		for x := 0; x < resp.Width; x++ {
			for y := 0; y < resp.Height; y++ {
				score := resp.At(x, y)
				if score < minscore {
					continue
				}
				pos := imgpyr.Point{k, image.Pt(x, y)}
				det := Det{score, pos}
				dets = append(dets, det)
			}
		}
	}
	log.Print("extracted windows: ", len(dets))
	return dets
}

// Extracts position and score per pixel in the response images.
func peaks(resps []*rimg64.Image, minscore float64) []Det {
	var dets []Det
	for k, resp := range resps {
		for x := 0; x < resp.Width; x++ {
			for y := 0; y < resp.Height; y++ {
				score := resp.At(x, y)
				if score < minscore {
					continue
				}
				if x > 0 && score < resp.At(x-1, y) {
					continue
				}
				if y > 0 && score < resp.At(x, y-1) {
					continue
				}
				if x < resp.Width-1 && score < resp.At(x+1, y) {
					continue
				}
				if y < resp.Height-1 && score < resp.At(x, y+1) {
					continue
				}
				pos := imgpyr.Point{k, image.Pt(x, y)}
				det := Det{score, pos}
				dets = append(dets, det)
			}
		}
	}
	log.Print("extracted peaks: ", len(dets))
	return dets
}

// Converts a point in the feature pyramid to a rectangle in the image.
func pointAt(pt imgpyr.Point, scales imgpyr.GeoSeq, rate int) image.Point {
	scale := scales.At(pt.Level)
	return vec(pt.Pos).Mul(float64(rate)).Mul(1 / scale).Round()
}

// Converts a point in the feature pyramid to a rectangle in the image.
func rectAt(pt imgpyr.Point, scales imgpyr.GeoSeq, rate int, interior image.Rectangle) image.Rectangle {
	scale := scales.At(pt.Level)
	a := vec(pt.Pos).Mul(float64(rate)).Add(vec(interior.Min)).Mul(1 / scale)
	// Scale position by rate, add size, scale by magnification.
	b := vec(pt.Pos).Mul(float64(rate)).Add(vec(interior.Max)).Mul(1 / scale)
	return image.Rectangle{a.Round(), b.Round()}
}

func intersect(a, b image.Rectangle, maxinter float64) bool {
	ab := a.Intersect(b)
	rela := float64(area(ab)) / float64(area(a))
	relb := float64(area(ab)) / float64(area(b))
	return (rela > maxinter && relb > maxinter)
}

func area(r image.Rectangle) int {
	s := r.Size()
	return s.X * s.Y
}

type byScoreDesc []Det

func (s byScoreDesc) Len() int           { return len(s) }
func (s byScoreDesc) Less(i, j int) bool { return s[i].Score > s[j].Score }
func (s byScoreDesc) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

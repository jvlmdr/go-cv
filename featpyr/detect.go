package featpyr

import (
	"github.com/jackvalmadre/go-cv/detect"
	"github.com/jackvalmadre/go-cv/imgpyr"
	"github.com/jackvalmadre/go-cv/rimg64"
	"github.com/jackvalmadre/go-cv/slide"

	"image"
	"log"
)

// Detection in pyramid.
type pyrdet struct {
	Score float64
	imgpyr.Point
}

type DetectOpts struct {
	// Limits the intersection of two detections.
	// The intersection is normalized by the union and is therefore between 0 and 1.
	// If MaxInter is zero then windows cannot overlap at all.
	// If MaxInter is one then windows can overlap arbitrarily.
	MaxInter float64
	// Threshold on scores. Can be negative infinity.
	MinScore float64
	// The number of detections.
	// If MaxNum <= 0 then the number of detections is unrestricted.
	MaxNum int
	// If LocalMax is true then every detection must be
	// greater than or equal to its 4-connected neighborhood.
	LocalMax bool
}

// Searches each level of the pyramid with a sliding window.
// Returns the detections ordered by score, from highest to lowest.
//
// The size of the template in pixels must be provided.
func Detect(pyr *Pyramid, tmpl *rimg64.Multi, pixsize image.Point, opts DetectOpts) []detect.Det {
	resps := evalTmpl(pyr, tmpl)
	pts := filterDets(resps, opts.MinScore, opts.LocalMax)
	rects := pointsToRects(pyr, pts, pixsize)
	return detect.SuppressOverlap(rects, opts.MaxNum, opts.MaxInter)
}

func pointsToRects(pyr *Pyramid, pts []pyrdet, pixsize image.Point) []detect.Det {
	dets := make([]detect.Det, len(pts))
	for i, pt := range pts {
		rect := rectAt(pt.Point, pyr.Images.Scales, pyr.Rate, pixsize)
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

func filterDets(resps []*rimg64.Image, minscore float64, localmax bool) []pyrdet {
	// Evaluate sliding window at each level.
	var dets []pyrdet
	if localmax {
		dets = peaks(resps, minscore)
	} else {
		// Put every window into a list.
		dets = windows(resps, minscore)
	}
	return dets
}

// Extracts position and score per pixel in the response images.
func windows(resps []*rimg64.Image, minscore float64) []pyrdet {
	var dets []pyrdet
	for k, resp := range resps {
		for x := 0; x < resp.Width; x++ {
			for y := 0; y < resp.Height; y++ {
				score := resp.At(x, y)
				if score < minscore {
					continue
				}
				pos := imgpyr.Point{k, image.Pt(x, y)}
				det := pyrdet{score, pos}
				dets = append(dets, det)
			}
		}
	}
	log.Print("extracted windows: ", len(dets))
	return dets
}

// Extracts position and score per pixel in the response images.
func peaks(resps []*rimg64.Image, minscore float64) []pyrdet {
	var dets []pyrdet
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
				det := pyrdet{score, pos}
				dets = append(dets, det)
			}
		}
	}
	log.Print("extracted peaks: ", len(dets))
	return dets
}

func (pyr *Pyramid) rectAt(pt imgpyr.Point, pixsize image.Point) image.Rectangle {
	return rectAt(pt, pyr.Images.Scales, pyr.Rate, pixsize)
}

// Converts a point in the feature pyramid to a rectangle in the image.
func rectAt(pt imgpyr.Point, scales imgpyr.GeoSeq, rate int, pixsize image.Point) image.Rectangle {
	scale := scales.At(pt.Level)
	a := vec(pt.Pos).Mul(float64(rate)).Mul(1 / scale)
	// Scale position by rate, add size, scale by magnification.
	b := vec(pt.Pos).Mul(float64(rate)).Add(vec(pixsize)).Mul(1 / scale)
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

type byScore []pyrdet

func (s byScore) Len() int           { return len(s) }
func (s byScore) Less(i, j int) bool { return s[i].Score < s[j].Score }
func (s byScore) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

package detect

import (
	"image"

	"github.com/jvlmdr/go-cv/feat"
	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-cv/slide"
)

// Detect performs detection and non-max suppression in a feature image.
// It returns a list of scored rectangles in the original image (before padding and feature transform).
// It is necessary to provide both the integer downsample rate of the feature transform
// and the margin which was added to the image before taking the feature transform.
// Calls detect.Score, detect.Sort, detect.Suppress.
func Detect(im *rimg64.Multi, margin feat.Margin, rate int, scorer slide.Scorer, shape PadRect, detopts DetFilter, suppropts SupprFilter) ([]Det, error) {
	dets, err := Score(im, margin, rate, scorer, shape, detopts)
	if err != nil {
		return nil, err
	}
	Sort(dets)
	dets = Suppress(dets, suppropts.MaxNum, suppropts.Overlap)
	return dets, nil
}

// DetectImage computes features and then calls Detect.
func DetectImage(im image.Image, phi feat.Image, pad feat.Pad, scorer slide.Scorer, shape PadRect, detopts DetFilter, suppropts SupprFilter) ([]Det, error) {
	f, err := feat.ApplyPad(phi, im, pad)
	if err != nil {
		return nil, err
	}
	return Detect(f, pad.Margin, phi.Rate(), scorer, shape, detopts, suppropts)
}

// Score computes the score of all windows in an image at a single scale.
// It does not perform non-max suppression.
// It returns an unordered list of scored rectangles.
func Score(im *rimg64.Multi, margin feat.Margin, rate int, scorer slide.Scorer, shape PadRect, opts DetFilter) ([]Det, error) {
	// Evaluate detector at all positions.
	pts, err := Points(im, scorer, opts.LocalMax, opts.MinScore)
	if err != nil {
		return nil, err
	}
	// Convert positions in the feature image to rectangles in the original image.
	dets := make([]Det, len(pts))
	for i, det := range pts {
		rect := featPtToImRect(det.Point, rate, margin, shape.Int)
		dets[i] = Det{det.Score, rect}
	}
	return dets, nil
}

// DetFilter specifies options to eliminate detections before non-max suppression.
type DetFilter struct {
	// Ignore detections which are smaller than a neighbor?
	LocalMax bool
	// Score threshold.
	MinScore float64
}

// DetPos describes a scored position.
type DetPos struct {
	Score float64
	image.Point
}

// Points performs detection without non-max suppression.
// It returns a list of unsorted scored positions in the feature image.
//
// If localmax is true, then points which have a neighbor greater than them are excluded.
// Any windows less than minscore are excluded.
func Points(im *rimg64.Multi, scorer slide.Scorer, localmax bool, minscore float64) ([]DetPos, error) {
	resp, err := slide.Score(im, scorer)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	var dets []DetPos
	// Iterate over positions and check criteria.
	for u := 0; u < resp.Width; u++ {
		for v := 0; v < resp.Height; v++ {
			score := resp.At(u, v)
			if score < minscore {
				continue
			}
			if localmax && notLocalMax(resp, u, v) {
				continue
			}
			// Scale by rate, then apply margin and interior offsets.
			dets = append(dets, DetPos{score, image.Pt(u, v)})
		}
	}
	return dets, nil
}

// Converts the position of a detection in a feature image to a rectangle in the intensity image.
//
// Additional arguments are:
// the integer downsample rate of the feature transform,
// the margin which was added to the image before taking the feature transform,
// the rectangular region within the window which corresponds to the annotation.
func featPtToImRect(pt image.Point, rate int, margin feat.Margin, interior image.Rectangle) image.Rectangle {
	return interior.Add(pt.Mul(rate)).Sub(margin.TopLeft())
}

// Tests whether (u, v) is a local maximum.
// Pixels at the edge can be maxima.
func notLocalMax(r *rimg64.Image, u, v int) bool {
	uv := r.At(u, v)
	if u > 0 && r.At(u-1, v) > uv {
		return true
	}
	if u < r.Width-1 && r.At(u+1, v) > uv {
		return true
	}
	if v > 0 && r.At(u, v-1) > uv {
		return true
	}
	if v < r.Height-1 && r.At(u, v+1) > uv {
		return true
	}
	return false
}

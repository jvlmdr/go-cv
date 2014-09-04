package detect

import (
	"image"

	"github.com/jvlmdr/go-cv/feat"
	"github.com/jvlmdr/go-cv/featpyr"
	"github.com/jvlmdr/go-cv/imgpyr"
	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/nfnt/resize"
)

// MultiScaleOpts specifies the parameters to MultiScale().
type MultiScaleOpts struct {
	PyrStep float64
	Interp  resize.InterpolationFunction
	feat.Transform
	feat.Pad
	DetFilter
	OverlapFunc
}

// MultiScale searches an image at multiple scales and performs non-max suppression.
//
// At each level, the image is rescaled using Interp,
// then padded using Pad before calling Transform.Apply().
// The levels are geometrically spaced at intervals of PyrStep.
// Detections are filtered using DetFtiler and then non-max suppression
// is performed using the OverlapFunc test.
func MultiScale(tmpl *FeatTmpl, im image.Image, opts MultiScaleOpts) []Det {
	scales := imgpyr.Scales(im.Bounds().Size(), tmpl.Size, opts.PyrStep)
	pyr := featpyr.NewPad(imgpyr.NewInterp(im, scales, opts.Interp), opts.Transform, opts.Pad)
	// Do not impose a maximum number of detections at each layer.
	supprOpts := SupprFilter{MaxNum: 0, Overlap: opts.OverlapFunc}
	return Pyramid(pyr, tmpl, opts.DetFilter, supprOpts)
}

// Performs detection and non-max suppression.
// Returns a list of scored detection windows.
// Windows are specified as rectangles in the original pixel image.
func Pyramid(pyr *featpyr.Pyramid, tmpl *FeatTmpl, detopts DetFilter, suppropts SupprFilter) []Det {
	// Get detections as top-left corners at some level.
	featdets := PyramidPoints(pyr, tmpl.Image, detopts.LocalMax, detopts.MinScore)
	// Convert to rectangles in the image.
	dets := make([]Det, len(featdets))
	for i, featdet := range featdets {
		rect := pyr.ToImageRect(featdet.Point, tmpl.Interior)
		dets[i] = Det{featdet.Score, rect}
	}
	// Non-max suppression.
	Sort(dets)
	dets = Suppress(dets, suppropts.MaxNum, suppropts.Overlap)
	return dets
}

// Scored position in feature pyramid.
type PyrDetPos struct {
	Score float64
	imgpyr.Point
}

// Returns scored windows in image.
// Windows are represented by the position of their top-left corner in the feature pyramid.
func PyramidPoints(pyr *featpyr.Pyramid, tmpl *rimg64.Multi, localmax bool, minscore float64) []PyrDetPos {
	var dets []PyrDetPos
	for level, im := range pyr.Feats {
		if im.Width < tmpl.Width || im.Height < tmpl.Height {
			break
		}
		// Get points from each level.
		imdets := Points(im, tmpl, localmax, minscore)
		// Append level to each point.
		for _, imdet := range imdets {
			pyrpt := imgpyr.Point{level, imdet.Point}
			dets = append(dets, PyrDetPos{imdet.Score, pyrpt})
		}
	}
	return dets
}

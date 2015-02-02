package detect

import (
	"image"
	"time"

	"github.com/jvlmdr/go-cv/feat"
	"github.com/jvlmdr/go-cv/featpyr"
	"github.com/jvlmdr/go-cv/imgpyr"
	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/nfnt/resize"
)

// MultiScaleOpts specifies the parameters to MultiScale().
type MultiScaleOpts struct {
	MaxScale  float64
	PyrStep   float64
	Interp    resize.InterpolationFunction
	Transform feat.Image
	feat.Pad
	DetFilter
	SupprFilter
}

type MultiScaleDuration struct{ Resize, Feat, Slide, Suppr time.Duration }

// MultiScale searches an image at multiple scales and performs non-max suppression.
//
// At each level, the image is rescaled using Interp,
// then padded using Pad before calling Transform.Apply().
// The levels are geometrically spaced at intervals of PyrStep.
// Detections are filtered using DetFilter and then non-max suppression
// is performed using the OverlapFunc test.
func MultiScale(im image.Image, tmpl *FeatTmpl, opts MultiScaleOpts) ([]Det, MultiScaleDuration, error) {
	scales := imgpyr.Scales(im.Bounds().Size(), tmpl.Size, opts.MaxScale, opts.PyrStep).Elems()
	ims := imgpyr.NewGenerator(im, scales, opts.Interp)
	pyr := featpyr.NewGenerator(ims, opts.Transform, opts.Pad)
	var dets []Det
	l, err := pyr.First()
	if err != nil {
		return nil, MultiScaleDuration{}, err
	}
	var dur MultiScaleDuration
	for l != nil {
		t := time.Now()
		pts := Points(l.Feat, tmpl.Image, tmpl.Bias, opts.DetFilter.LocalMax, opts.DetFilter.MinScore)
		dur.Slide += time.Since(t)
		// Convert to scored rectangles in the image.
		for _, pt := range pts {
			rect := pyr.ToImageRect(l.Image.Index, pt.Point, tmpl.Interior)
			dets = append(dets, Det{pt.Score, rect})
		}
		var err error
		l, err = pyr.Next(l)
		if err != nil {
			return nil, MultiScaleDuration{}, err
		}
	}
	dur.Resize = pyr.DurResize
	dur.Feat = pyr.DurFeat
	t := time.Now()
	Sort(dets)
	dets = Suppress(dets, opts.SupprFilter.MaxNum, opts.SupprFilter.Overlap)
	dur.Suppr = time.Since(t)
	return dets, dur, nil
}

// Performs detection and non-max suppression.
// Returns a list of scored detection windows.
// Windows are specified as rectangles in the original pixel image.
func Pyramid(pyr *featpyr.Pyramid, tmpl *FeatTmpl, detopts DetFilter, suppropts SupprFilter) []Det {
	// Get detections as top-left corners at some level.
	featdets := detectPyrPoints(pyr, tmpl.Image, tmpl.Bias, detopts.LocalMax, detopts.MinScore)
	// Convert to rectangles in the image.
	dets := make([]Det, len(featdets))
	for i, featdet := range featdets {
		rect := pyr.ToImageRect(featdet.Point, tmpl.Interior)
		dets[i] = Det{featdet.Score, rect}
	}
	// Non-max suppression.
	Sort(dets)
	return Suppress(dets, suppropts.MaxNum, suppropts.Overlap)
}

// Scored position in feature pyramid.
type pyrDetPos struct {
	Score float64
	imgpyr.Point
}

// Returns scored windows in image.
// Windows are represented by the position of their top-left corner in the feature pyramid.
func detectPyrPoints(pyr *featpyr.Pyramid, tmpl *rimg64.Multi, bias float64, localmax bool, minscore float64) []pyrDetPos {
	var dets []pyrDetPos
	for level, im := range pyr.Feats {
		if im.Width < tmpl.Width || im.Height < tmpl.Height {
			break
		}
		// Get points from each level.
		imdets := Points(im, tmpl, bias, localmax, minscore)
		// Append level to each point.
		for _, imdet := range imdets {
			pyrpt := imgpyr.Point{level, imdet.Point}
			dets = append(dets, pyrDetPos{imdet.Score, pyrpt})
		}
	}
	return dets
}

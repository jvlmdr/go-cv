package detect

import (
	"image"
	"time"

	"github.com/jvlmdr/go-cv/feat"
	"github.com/jvlmdr/go-cv/featpyr"
	"github.com/jvlmdr/go-cv/imgpyr"
	"github.com/jvlmdr/go-cv/slide"
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
func MultiScale(im image.Image, scorer slide.Scorer, shape PadRect, opts MultiScaleOpts) ([]Det, MultiScaleDuration, error) {
	scales := imgpyr.Scales(im.Bounds().Size(), scorer.Size(), opts.MaxScale, opts.PyrStep).Elems()
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
		pts, err := Points(l.Feat, scorer, opts.DetFilter.LocalMax, opts.DetFilter.MinScore)
		if err != nil {
			return nil, MultiScaleDuration{}, err
		}
		dur.Slide += time.Since(t)
		// Convert to scored rectangles in the image.
		for _, pt := range pts {
			rect := pyr.ToImageRect(l.Image.Index, pt.Point, shape.Int)
			dets = append(dets, Det{pt.Score, rect})
		}
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

// Pyramid performs detection and non-max suppression.
// Returns a list of scored detection windows.
// Windows are specified as rectangles in the original pixel image.
func Pyramid(pyr *featpyr.Pyramid, scorer slide.Scorer, shape PadRect, detopts DetFilter, suppropts SupprFilter) ([]Det, error) {
	// Get detections as top-left corners at some level.
	featdets, err := detectPyrPoints(pyr, scorer, detopts.LocalMax, detopts.MinScore)
	if err != nil {
		return nil, err
	}
	// Convert to rectangles in the image.
	dets := make([]Det, len(featdets))
	for i, featdet := range featdets {
		rect := pyr.ToImageRect(featdet.Point, shape.Int)
		dets[i] = Det{featdet.Score, rect}
	}
	// Non-max suppression.
	Sort(dets)
	return Suppress(dets, suppropts.MaxNum, suppropts.Overlap), nil
}

// Scored position in feature pyramid.
type pyrDetPos struct {
	Score float64
	imgpyr.Point
}

// Returns scored windows in image.
// Windows are represented by the position of their top-left corner in the feature pyramid.
func detectPyrPoints(pyr *featpyr.Pyramid, scorer slide.Scorer, localmax bool, minscore float64) ([]pyrDetPos, error) {
	var dets []pyrDetPos
	size := scorer.Size()
	for level, im := range pyr.Feats {
		if im.Width < size.X || im.Height < size.Y {
			break
		}
		// Get points from each level.
		imdets, err := Points(im, scorer, localmax, minscore)
		if err != nil {
			return nil, err
		}
		// Append level to each point.
		for _, imdet := range imdets {
			pyrpt := imgpyr.Point{level, imdet.Point}
			dets = append(dets, pyrDetPos{imdet.Score, pyrpt})
		}
	}
	return dets, nil
}

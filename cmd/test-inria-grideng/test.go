package main

import (
	"image"
	"log"
	"math"
	"path"

	"github.com/jackvalmadre/go-cv/dataset/inria"
	"github.com/jackvalmadre/go-cv/detect"
	"github.com/jackvalmadre/go-cv/feat"
	"github.com/jackvalmadre/go-cv/featpyr"
	"github.com/jackvalmadre/go-cv/hog"
	"github.com/jackvalmadre/go-cv/imgpyr"
	"github.com/jackvalmadre/go-cv/imsamp"
	"github.com/jackvalmadre/go-grideng/grideng"
)

type ValidateArgs struct {
	Tmpl *detect.FeatTmpl
	Dir  string
	DetectOpts
	MinInter float64
}

func init() {
	grideng.Register("validate", grideng.Func(
		func(annot inria.Annot, p ValidateArgs) (*detect.ResultSet, error) {
			return testImage(p.Tmpl, annot, p.Dir, p.DetectOpts, p.MinInter)
		},
	))
}

type DetectOpts struct {
	HOGBin   int
	PyrStep  float64
	MaxIOU   float64
	Margin   int
	LocalMax bool
}

type OverlapTest struct {
	// Overlap if Cover(a, b) >= MaxCover.
	MaxCover float64
	// Or if IOU(a, b) >= MaxIOU.
	MaxIOU float64
	// Or if Cover(a, b) >= MaxCoverBoth && Cover(b, a) >= MaxCoverBoth.
	MaxCoverBoth float64
}

func (t OverlapTest) Overlap(a, b image.Rectangle) bool {
	if detect.Cover(a, b) >= t.MaxCover {
		return true
	}
	if detect.IOU(a, b) >= t.MaxIOU {
		return true
	}
	if detect.Cover(a, b) >= t.MaxCoverBoth && detect.Cover(b, a) >= t.MaxCoverBoth {
		return true
	}
	return false
}

func (t OverlapTest) Func() detect.OverlapFunc {
	return func(a, b image.Rectangle) bool { return t.Overlap(a, b) }
}

func test(tmpl *detect.FeatTmpl, annots []inria.Annot, dir string, opts DetectOpts, mininter float64) (*detect.ResultSet, error) {
	// Execute in parallel.
	vals := make([]*detect.ResultSet, len(annots))
	conf := ValidateArgs{
		Tmpl:       tmpl,
		Dir:        dir,
		DetectOpts: opts,
		MinInter:   mininter,
	}
	if err := grideng.Map("validate", vals, annots, conf); err != nil {
		log.Fatalln("validate:", err)
	}
	return detect.MergeResults(vals...), nil
}

// Runs detector across a single image and validates results.
func testImage(tmpl *detect.FeatTmpl, annot inria.Annot, dir string, opts DetectOpts, mininter float64) (*detect.ResultSet, error) {
	im, err := loadImage(path.Join(dir, annot.Image))
	if err != nil {
		return nil, err
	}
	// Get detections.
	dets := detectImage(tmpl, im, opts.Margin, opts.PyrStep, opts.HOGBin, opts.LocalMax, opts.MaxIOU)
	val := detect.ValidateMatch(dets, annot.Rects, mininter)
	return val, nil
}

// Runs a single detector across a single image and returns results.
func detectImage(tmpl *detect.FeatTmpl, im image.Image, margin int, step float64, sbin int, localmax bool, maxiou float64) []detect.Det {
	// Construct pyramid.
	// Get range of scales.
	scales := imgpyr.Scales(im.Bounds().Size(), tmpl.Size, step)
	// Define feature transform.
	phi := hog.Transform{hog.FGMRConfig(sbin)}
	// Define amount and type of padding.
	pad := feat.Pad{feat.Margin{margin, margin, margin, margin}, imsamp.Continue}
	pyr := featpyr.NewPad(imgpyr.New(im, scales), phi, pad)

	// Search feature pyramid.
	// Options for running detector on each level.
	detopts := detect.DetFilter{LocalMax: localmax, MinScore: math.Inf(-1)}
	// Use intersection-over-union criteria for non-max suppression.
	overlap := func(a, b image.Rectangle) bool {
		return detect.IOU(a, b) > maxiou
	}
	// Options for non-max suppression.
	suppropts := detect.SupprFilter{MaxNum: 0, Overlap: overlap}
	dets := detect.Pyramid(pyr, tmpl, detopts, suppropts)
	return dets
}

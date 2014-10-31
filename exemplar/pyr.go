package exemplar

import (
	"image"

	"github.com/jvlmdr/go-cv/detect"
	"github.com/jvlmdr/go-cv/featpyr"
	"github.com/jvlmdr/go-cv/imgpyr"
)

func MultiScale(im image.Image, tmpls map[string]*detect.FeatTmpl, opts detect.MultiScaleOpts) ([]Det, error) {
	if len(tmpls) == 0 {
		return nil, nil
	}
	scales := imgpyr.Scales(im.Bounds().Size(), minDims(tmpls), opts.MaxScale, opts.PyrStep).Elems()
	ims := imgpyr.NewGenerator(im, scales, opts.Interp)
	pyr := featpyr.NewGenerator(ims, opts.Transform, opts.Pad)
	var dets []Det
	l, err := pyr.First()
	if err != nil {
		return nil, err
	}
	for l != nil {
		for key, tmpl := range tmpls {
			pts := detect.Points(l.Feat, tmpl.Image, opts.DetFilter.LocalMax, opts.DetFilter.MinScore)
			// Convert to scored rectangles in the image.
			for _, pt := range pts {
				rect := pyr.ToImageRect(l.Image.Index, pt.Point, tmpl.Interior)
				dets = append(dets, Det{detect.Det{pt.Score + tmpl.Bias, rect}, key})
			}
		}
		var err error
		l, err = pyr.Next(l)
		if err != nil {
			return nil, err
		}
	}
	Sort(dets)
	inds := detect.SuppressIndex(DetSlice(dets), opts.SupprFilter.MaxNum, opts.SupprFilter.Overlap)
	dets = detsSubset(dets, inds)
	return dets, nil
}

func detsSubset(dets []Det, inds []int) []Det {
	subset := make([]Det, len(inds))
	for p, i := range inds {
		subset[p] = dets[i]
	}
	return subset
}

func minDims(tmpls map[string]*detect.FeatTmpl) image.Point {
	var (
		x, y int
		init bool
	)
	for _, tmpl := range tmpls {
		if !init {
			x, y, init = tmpl.Size.X, tmpl.Size.Y, true
			continue
		}
		if tmpl.Size.X < x {
			x = tmpl.Size.X
		}
		if tmpl.Size.Y < y {
			y = tmpl.Size.Y
		}
	}
	return image.Pt(x, y)
}

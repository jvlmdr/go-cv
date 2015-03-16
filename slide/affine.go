package slide

import (
	"fmt"
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
)

type CorrOp int

const (
	Dot CorrOp = iota
	Cos
)

type AffineScorer struct {
	Tmpl *rimg64.Multi
	Bias float64
	Op   CorrOp
}

func (f *AffineScorer) Size() image.Point {
	return f.Tmpl.Size()
}

func (f *AffineScorer) Score(x *rimg64.Multi) (float64, error) {
	if f.Op != Cos {
		panic("cosine unimplemented")
	}
	if !x.Size().Eq(f.Tmpl.Size()) {
		return 0, fmt.Errorf("different size: input %v, template %v", x.Size(), f.Tmpl.Size())
	}
	if x.Channels != f.Tmpl.Channels {
		return 0, fmt.Errorf("different channels: input %v, template %v", x.Channels, f.Tmpl.Channels)
	}
	size := f.Tmpl.Size()
	var y float64
	for i := 0; i < size.X; i++ {
		for j := 0; j < size.Y; j++ {
			for k := 0; k < f.Tmpl.Channels; k++ {
				y += x.At(i, j, k) * f.Tmpl.At(i, j, k)
			}
		}
	}
	y += f.Bias
	return y, nil
}

func (f *AffineScorer) Slide(im *rimg64.Multi) (*rimg64.Image, error) {
	if f.Op != Dot {
		panic("unimplemented")
	}
	y, err := CorrMultiAuto(im, f.Tmpl)
	if err != nil {
		return nil, err
	}
	if f.Bias == 0 {
		return y, nil
	}
	for i := 0; i < y.Width; i++ {
		for j := 0; j < y.Height; j++ {
			y.Set(i, j, y.At(i, j)+f.Bias)
		}
	}
	return y, nil
}

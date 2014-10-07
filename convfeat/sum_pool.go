package convfeat

import (
	"fmt"
	"image"

	"github.com/jvlmdr/go-cv/feat"
	"github.com/jvlmdr/go-cv/rimg64"
)

func init() {
	feat.RegisterReal("sum-pool", func() feat.Real { return new(SumPool) })
}

type SumPool struct {
	Field  image.Point
	Stride int
}

func (phi SumPool) Rate() int { return phi.Stride }

func (phi SumPool) Apply(x *rimg64.Multi) (*rimg64.Multi, error) {
	if phi.Field.X <= 0 || phi.Field.Y <= 0 {
		err := fmt.Errorf("invalid field size: %v", phi.Field)
		return nil, err
	}
	if phi.Stride <= 0 {
		err := fmt.Errorf("invalid stride: %d", phi.Stride)
		return nil, err
	}
	size := x.Size().Sub(phi.Field).Add(image.Pt(1, 1)).Div(phi.Stride)
	y := rimg64.NewMulti(size.X, size.Y, x.Channels)
	for i := 0; i < y.Width; i++ {
		for j := 0; j < y.Height; j++ {
			for k := 0; k < x.Channels; k++ {
				// Position in original image.
				p := image.Pt(i, j).Mul(phi.Stride)
				var t float64
				for u := p.X; u < p.X+phi.Field.X; u++ {
					for v := p.Y; v < p.Y+phi.Field.Y; v++ {
						t += x.At(u, v, k)
					}
				}
				y.Set(i, j, k, t)
			}
		}
	}
	return y, nil
}

package convfeat

import (
	"fmt"
	"image"
	"math"

	"github.com/jvlmdr/go-cv/feat"
	"github.com/jvlmdr/go-cv/rimg64"
)

func init() {
	feat.RegisterReal("max-pool", func() feat.RealSpec {
		return feat.NewRealSpec(new(MaxPool))
	})
}

type MaxPool struct {
	Field  image.Point
	Stride int
}

func (phi *MaxPool) Rate() int { return phi.Stride }

func (phi *MaxPool) Apply(x *rimg64.Multi) (*rimg64.Multi, error) {
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
				max := math.Inf(-1)
				for u := 0; u < phi.Field.X; u++ {
					for v := 0; v < phi.Field.Y; v++ {
						q := p.Add(image.Pt(u, v))
						max = math.Max(max, x.At(q.X, q.Y, k))
					}
				}
				y.Set(i, j, k, max)
			}
		}
	}
	return y, nil
}

func (phi *MaxPool) Marshaler() *feat.RealMarshaler {
	return &feat.RealMarshaler{"max-pool", feat.NewRealSpec(phi)}
}

package convfeat

import (
	"fmt"
	"image"
	"math"

	"github.com/jvlmdr/go-cv/featset"
	"github.com/jvlmdr/go-cv/rimg64"
)

func init() {
	featset.RegisterReal("max-pool", func() featset.Real { return new(MaxPool) })
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
	size := image.Pt(
		ceilDiv(x.Width-phi.Field.X+1, phi.Stride),
		ceilDiv(x.Height-phi.Field.Y+1, phi.Stride),
	)
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

func (phi *MaxPool) Marshaler() *featset.RealMarshaler {
	return &featset.RealMarshaler{"max-pool", phi}
}

func (phi *MaxPool) Transform() featset.Real { return phi }

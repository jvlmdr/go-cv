package convfeat

import (
	"math"

	"github.com/jvlmdr/go-cv/featset"
	"github.com/jvlmdr/go-cv/rimg64"
)

func init() {
	featset.RegisterReal("pos", func() featset.Real { return new(PosPart) })
	featset.RegisterReal("pos-neg", func() featset.Real { return new(PosNegPart) })
	featset.RegisterReal("is-pos", func() featset.Real { return new(IsPos) })
	featset.RegisterReal("sign", func() featset.Real { return new(Sign) })
}

// PosPart takes the positive part of the input.
type PosPart struct{}

func (phi *PosPart) Rate() int { return 1 }

func (phi *PosPart) Apply(x *rimg64.Multi) (*rimg64.Multi, error) {
	y := rimg64.NewMulti(x.Width, x.Height, x.Channels)
	for i := 0; i < x.Width; i++ {
		for j := 0; j < x.Height; j++ {
			for k := 0; k < x.Channels; k++ {
				y.Set(i, j, k, math.Abs(x.At(i, j, k)))
			}
		}
	}
	return y, nil
}

func (phi *PosPart) Marshaler() *featset.RealMarshaler {
	return &featset.RealMarshaler{"pos", nil}
}

func (phi *PosPart) Transform() featset.Real { return phi }

// PosNegPart splits the input into its positive and negative parts.
type PosNegPart struct{}

func (phi *PosNegPart) Rate() int { return 1 }

func (phi *PosNegPart) Apply(x *rimg64.Multi) (*rimg64.Multi, error) {
	channels := x.Channels * 2
	y := rimg64.NewMulti(x.Width, x.Height, channels)
	for i := 0; i < x.Width; i++ {
		for j := 0; j < x.Height; j++ {
			for k := 0; k < x.Channels; k++ {
				pos, neg := posNegPart(x.At(i, j, k))
				y.Set(i, j, 2*k, pos)
				y.Set(i, j, 2*k+1, neg)
			}
		}
	}
	return y, nil
}

func (phi *PosNegPart) Marshaler() *featset.RealMarshaler {
	return &featset.RealMarshaler{"pos-neg", nil}
}

func (phi *PosNegPart) Transform() featset.Real { return phi }

func posNegPart(x float64) (pos, neg float64) {
	if math.IsNaN(x) {
		panic("nan")
	}
	if x < 0 {
		return 0, -x
	}
	return x, 0
}

// IsPos returns 1 if positive, 0 otherwise.
type IsPos struct{}

func (phi *IsPos) Rate() int { return 1 }

func (phi *IsPos) Apply(x *rimg64.Multi) (*rimg64.Multi, error) {
	y := rimg64.NewMulti(x.Width, x.Height, x.Channels)
	for i := 0; i < x.Width; i++ {
		for j := 0; j < x.Height; j++ {
			for k := 0; k < x.Channels; k++ {
				if x.At(i, j, k) > 0 {
					y.Set(i, j, k, 1)
				}
			}
		}
	}
	return y, nil
}

func (phi *IsPos) Marshaler() *featset.RealMarshaler {
	return &featset.RealMarshaler{"is-pos", nil}
}

func (phi *IsPos) Transform() featset.Real { return phi }

// Sign returns 1 if positive, -1 if negative.
type Sign struct{}

func (phi *Sign) Rate() int { return 1 }

func (phi *Sign) Apply(x *rimg64.Multi) (*rimg64.Multi, error) {
	y := rimg64.NewMulti(x.Width, x.Height, x.Channels)
	for i := 0; i < x.Width; i++ {
		for j := 0; j < x.Height; j++ {
			for k := 0; k < x.Channels; k++ {
				y.Set(i, j, k, sign(x.At(i, j, k)))
			}
		}
	}
	return y, nil
}

func (phi *Sign) Marshaler() *featset.RealMarshaler {
	return &featset.RealMarshaler{"sign", nil}
}

func (phi *Sign) Transform() featset.Real { return phi }

func sign(x float64) float64 {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return x
	}
}

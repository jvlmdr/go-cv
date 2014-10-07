package convfeat

import (
	"math"

	"github.com/jvlmdr/go-cv/feat"
	"github.com/jvlmdr/go-cv/rimg64"
)

func init() {
	feat.RegisterReal("pos", func() feat.Real { return new(PosPart) })
	feat.RegisterReal("pos-neg", func() feat.Real { return new(PosNegPart) })
	feat.RegisterReal("is-pos", func() feat.Real { return new(IsPos) })
	feat.RegisterReal("sign", func() feat.Real { return new(Sign) })
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

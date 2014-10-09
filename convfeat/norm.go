package convfeat

import (
	"math"

	"github.com/jvlmdr/go-cv/feat"
	"github.com/jvlmdr/go-cv/rimg64"
)

func init() {
	feat.RegisterReal("adj-chan-norm", func() feat.RealSpec {
		return feat.NewRealSpec(new(AdjChanNorm))
	})
}

// AdjChanNorm describes normalization over adjacent channels.
// For one vector-valued pixel a, it computes
// 	a[i] / (K + Alpha sum_{j = i-Num/2}^{i+Num/2} a[j]^2)^Beta
type AdjChanNorm struct {
	// Number of channels over which to normalize.
	// Must be odd.
	Num int
	// Parameters of normalization.
	K, Alpha, Beta float64
}

func (phi *AdjChanNorm) Rate() int { return 1 }

func (phi *AdjChanNorm) Apply(x *rimg64.Multi) (*rimg64.Multi, error) {
	y := rimg64.NewMulti(x.Width, x.Height, x.Channels)
	r := (phi.Num - 1) / 2
	for i := 0; i < x.Width; i++ {
		for j := 0; j < x.Height; j++ {
			k := 0
			// Range over which to compute sum.
			a, b := k-r, k+r+1
			// Take sum excluding leading element.
			var t float64
			for p := 0; p < min(b, x.Channels); p++ {
				t += sqr(x.At(i, j, p))
			}
			for ; k < x.Channels; k++ {
				a, b = k-r, k+r+1
				// Set element.
				norm := math.Pow(phi.K+phi.Alpha*t, phi.Beta)
				y.Set(i, j, k, x.At(i, j, k)/norm)
				// Subtract trailing element.
				if a >= 0 {
					t -= sqr(x.At(i, j, a))
				}
				// Add leading element.
				if b < x.Channels {
					t += sqr(x.At(i, j, b))
				}
			}
		}
	}
	return y, nil
}

func (phi *AdjChanNorm) Marshaler() *feat.RealMarshaler {
	return &feat.RealMarshaler{"adj-chan-norm", feat.NewRealSpec(phi)}
}

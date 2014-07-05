package detect

import (
	"math"
	"testing"
)

func TestSetAspect(t *testing.T) {
	const eps = 1e-9
	tests := []struct {
		InX, InY   float64
		Aspect     float64
		Mode       string
		OutX, OutY float64
	}{
		{1, 2, 1, "area", math.Sqrt(2), math.Sqrt(2)},
		{1, 1, 2, "area", math.Sqrt(2), 1 / math.Sqrt(2)},
		{1, 1, 0.5, "area", 1 / math.Sqrt(2), math.Sqrt(2)},
	}

	for _, e := range tests {
		x, y := SetAspect(e.InX, e.InY, e.Aspect, e.Mode)
		if math.Abs(x-e.OutX) > eps || math.Abs(y-e.OutY) > eps {
			t.Errorf(
				"{(%.4g, %.4g) Apsect:%.4g Mode:%s}: want (%.4g, %.4g), got (%g, %g)",
				e.InX, e.InY, e.Aspect, e.Mode, e.OutX, e.OutY, x, y,
			)
		}
	}
}

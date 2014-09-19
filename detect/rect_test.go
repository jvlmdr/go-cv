package detect_test

import (
	"math"
	"testing"

	"github.com/jvlmdr/go-cv/detect"
)

func TestSetAspect(t *testing.T) {
	const eps = 1e-9
	tests := []struct {
		InX, InY   float64
		Aspect     float64
		Mode       string
		OutX, OutY float64
	}{
		// Resize (1, 2) such that aspect ratio is 1.
		{1, 2, 1, "area", math.Sqrt(2), math.Sqrt(2)},
		{1, 2, 1, "width", 1, 1},
		{1, 2, 1, "height", 2, 2},
		{1, 2, 1, "fit", 1, 1},
		{1, 2, 1, "fill", 2, 2},

		// Resize (1, 1) such that aspect ratio (width / height) is 2.
		{1, 1, 2, "area", math.Sqrt(2), 1 / math.Sqrt(2)},
		{1, 1, 2, "width", 1, 0.5},
		{1, 1, 2, "height", 2, 1},
		{1, 1, 2, "fit", 1, 0.5},
		{1, 1, 2, "fill", 2, 1},

		// Resize (1, 1) such that aspect ratio (width / height) is 0.5.
		{1, 1, 0.5, "area", 1 / math.Sqrt(2), math.Sqrt(2)},
		{1, 1, 0.5, "width", 1, 2},
		{1, 1, 0.5, "height", 0.5, 1},
		{1, 1, 0.5, "fit", 0.5, 1},
		{1, 1, 0.5, "fill", 1, 2},

		// Resize (1, 2) such that aspect ratio (width / height) is 2.
		{1, 2, 2, "area", 2, 1},
		{1, 2, 2, "width", 1, 0.5},
		{1, 2, 2, "height", 4, 2},
		{1, 2, 2, "fit", 1, 0.5},
		{1, 2, 2, "fill", 4, 2},
	}

	for _, e := range tests {
		x, y := detect.SetAspect(e.InX, e.InY, e.Aspect, e.Mode)
		if math.Abs(x-e.OutX) > eps || math.Abs(y-e.OutY) > eps {
			t.Errorf(
				"{(%.4g, %.4g) Apsect:%.4g Mode:%s}: want (%.4g, %.4g), got (%g, %g)",
				e.InX, e.InY, e.Aspect, e.Mode, e.OutX, e.OutY, x, y,
			)
		}
	}
}

package detect_test

import (
	"image"
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

func TestFitRect(t *testing.T) {
	tests := []struct {
		Orig   image.Rectangle
		Target detect.PadRect
		Mode   string
		Out    image.Rectangle
	}{
		{
			// 20x40 box at position (30, 40).
			Orig: image.Rect(0, 0, 20, 40).Add(image.Pt(30, 40)),
			// 10x20 box at position (10, 10) in a 30x40 box.
			Target: detect.PadRect{
				Size: image.Pt(30, 40),
				Int:  image.Rect(0, 0, 10, 20).Add(image.Pt(10, 10)),
			},
			Mode: "area",
			Out:  image.Rect(0, 0, 20, 40).Add(image.Pt(30, 40)).Inset(-20),
		},
		{
			// 20x40 box at position (30, 40).
			Orig: image.Rect(0, 0, 20, 40).Add(image.Pt(30, 40)),
			// 10x20 box at position (5, 10) in a 30x40 box.
			// Therefore (15, 10) on other side.
			Target: detect.PadRect{
				Size: image.Pt(30, 40),
				Int:  image.Rect(0, 0, 10, 20).Add(image.Pt(5, 10)),
			},
			Mode: "area",
			Out: image.Rectangle{
				image.Pt(0, 0).Mul(2).Add(image.Pt(30, 40)).Sub(image.Pt(5, 10).Mul(2)),
				image.Pt(10, 20).Mul(2).Add(image.Pt(30, 40)).Add(image.Pt(15, 10).Mul(2)),
			},
		},
		{
			// 20x40 box at position (30, 40).
			Orig: image.Rect(0, 0, 20, 40).Add(image.Pt(30, 40)),
			// 10x10 box at position (10, 10) in a 30x40 box.
			// Therefore (10, 20) on other side.
			Target: detect.PadRect{
				Size: image.Pt(30, 40),
				Int:  image.Rect(0, 0, 10, 10).Add(image.Pt(10, 10)),
			},
			Mode: "height",
			// The centroid of the original rectangle is (30+10, 40+20) = (40, 60).
			// The centroid of the target's interior rectangle is (10+5, 10+5) = (15, 15).
			// From the centroid to the bottom right corner is (30-15, 40-15) = (15, 25).
			// The scale factor for the height is 4.
			// Thus the overall rectangle is (40-15*4, 60-15*4)-(40+15*4, 60+25*4) = (-20,0)-(100,160).
			// Its interior is (-20+10*4,0+10*4)-(100-10*4,160-20*4) = (20,40)-(60,80).
			Out: image.Rect(-20, 0, 100, 160),
		},
	}

	for _, e := range tests {
		_, got := detect.FitRect(e.Orig, e.Target, e.Mode)
		if !got.Eq(e.Out) {
			t.Errorf("orig %v, target %v, mode %s: want %v, got %v", e.Orig, e.Target, e.Mode, e.Out, got)
		}
	}
}

package imsamp

import (
	"image"
	"image/color"
	"math/rand"
	"testing"
)

func TestBlack(t *testing.T) {
	im := &image.Gray{
		[]uint8{
			10, 20, 30,
			40, 50, 60,
		},
		3,
		image.Rect(0, 0, 3, 2),
	}

	cases := []struct {
		Point image.Point
		Color color.Color
	}{
		{image.Pt(0, 0), color.Gray{10}},
		{image.Pt(2, 0), color.Gray{30}},
		{image.Pt(3, 0), color.Black},
		{image.Pt(2, 1), color.Gray{60}},
		{image.Pt(2, 2), color.Black},
		{image.Pt(3, 1), color.Black},
		{image.Pt(-1, 0), color.Black},
		{image.Pt(0, -1), color.Black},
		{image.Pt(0, 1), color.Gray{40}},
		{image.Pt(0, 2), color.Black},
		{image.Pt(-1, 2), color.Black},
		{image.Pt(-1, -1), color.Black},
	}

	for _, c := range cases {
		color := Black(im, c.Point)
		if !eq(color, c.Color) {
			t.Fatalf("at %v: want %v, got %v", c.Point, c.Color, color)
		}
	}
}

func TestPeriodic(t *testing.T) {
	im := &image.Gray{
		[]uint8{
			10, 20, 30,
			40, 50, 60,
		},
		3,
		image.Rect(0, 0, 3, 2),
	}

	cases := []struct {
		Point image.Point
		Color color.Color
	}{
		{image.Pt(0, 0), color.Gray{10}},
		{image.Pt(2, 0), color.Gray{30}},
		{image.Pt(3, 0), color.Gray{10}},
		{image.Pt(2, 1), color.Gray{60}},
		{image.Pt(2, 2), color.Gray{30}},
		{image.Pt(3, 1), color.Gray{40}},
		{image.Pt(-1, 0), color.Gray{30}},
		{image.Pt(0, -1), color.Gray{40}},
		{image.Pt(0, 1), color.Gray{40}},
		{image.Pt(0, 2), color.Gray{10}},
		{image.Pt(-1, 2), color.Gray{30}},
		{image.Pt(-1, -1), color.Gray{60}},
	}

	for _, c := range cases {
		color := Periodic(im, c.Point)
		if !eq(color, c.Color) {
			t.Fatalf("at %v: want %v, got %v", c.Point, c.Color, color)
		}
	}
}

func TestSymmetric(t *testing.T) {
	im := &image.Gray{
		[]uint8{
			10, 20, 30,
			40, 50, 60,
		},
		3,
		image.Rect(0, 0, 3, 2),
	}

	cases := []struct {
		Point image.Point
		Color color.Color
	}{
		{image.Pt(0, 0), color.Gray{10}},
		{image.Pt(2, 0), color.Gray{30}},
		{image.Pt(3, 0), color.Gray{30}},
		{image.Pt(4, 0), color.Gray{20}},
		{image.Pt(5, 0), color.Gray{10}},
		{image.Pt(6, 0), color.Gray{10}},
		{image.Pt(7, 0), color.Gray{20}},

		{image.Pt(2, 1), color.Gray{60}},
		{image.Pt(2, 2), color.Gray{60}},
		{image.Pt(2, 3), color.Gray{30}},
		{image.Pt(2, 4), color.Gray{30}},

		{image.Pt(3, 1), color.Gray{60}},

		{image.Pt(-1, 0), color.Gray{10}},
		{image.Pt(-2, 0), color.Gray{20}},
		{image.Pt(-3, 0), color.Gray{30}},
		{image.Pt(-4, 0), color.Gray{30}},
		{image.Pt(-5, 0), color.Gray{20}},

		{image.Pt(0, -1), color.Gray{10}},
		{image.Pt(0, -2), color.Gray{40}},
		{image.Pt(0, -3), color.Gray{40}},
		{image.Pt(0, -4), color.Gray{10}},

		{image.Pt(0, 1), color.Gray{40}},
		{image.Pt(0, 2), color.Gray{40}},
		{image.Pt(0, 3), color.Gray{10}},

		{image.Pt(-1, -1), color.Gray{10}},
		{image.Pt(-2, -2), color.Gray{50}},
		{image.Pt(-3, -3), color.Gray{60}},
		{image.Pt(-4, -4), color.Gray{30}},
		{image.Pt(-5, -5), color.Gray{20}},
	}

	for _, c := range cases {
		color := Symmetric(im, c.Point)
		if !eq(color, c.Color) {
			t.Fatalf("at %v: want %v, got %v", c.Point, c.Color, color)
		}
	}
}

func TestContinue(t *testing.T) {
	im := &image.Gray{
		[]uint8{
			10, 20, 30,
			40, 50, 60,
		},
		3,
		image.Rect(0, 0, 3, 2),
	}

	cases := []struct {
		Point image.Point
		Color color.Color
	}{
		{image.Pt(0, 0), color.Gray{10}},
		{image.Pt(2, 0), color.Gray{30}},
		{image.Pt(3, 0), color.Gray{30}},
		{image.Pt(4, 0), color.Gray{30}},
		{image.Pt(5, 0), color.Gray{30}},

		{image.Pt(2, 1), color.Gray{60}},
		{image.Pt(2, 2), color.Gray{60}},
		{image.Pt(2, 3), color.Gray{60}},
		{image.Pt(2, 4), color.Gray{60}},

		{image.Pt(3, 1), color.Gray{60}},

		{image.Pt(-1, 0), color.Gray{10}},
		{image.Pt(-2, 0), color.Gray{10}},
		{image.Pt(-3, 0), color.Gray{10}},
		{image.Pt(-4, 0), color.Gray{10}},

		{image.Pt(0, -1), color.Gray{10}},
		{image.Pt(0, -2), color.Gray{10}},
		{image.Pt(0, -3), color.Gray{10}},
		{image.Pt(0, -4), color.Gray{10}},

		{image.Pt(0, 1), color.Gray{40}},
		{image.Pt(0, 2), color.Gray{40}},
		{image.Pt(0, 3), color.Gray{40}},

		{image.Pt(-1, -1), color.Gray{10}},
		{image.Pt(-2, -2), color.Gray{10}},
		{image.Pt(-3, -3), color.Gray{10}},
		{image.Pt(-4, -4), color.Gray{10}},
	}

	for _, c := range cases {
		color := Continue(im, c.Point)
		if !eq(color, c.Color) {
			t.Fatalf("at %v: want %v, got %v", c.Point, c.Color, color)
		}
	}
}

func eq(u, v color.Color) bool {
	ru, gu, bu, au := u.RGBA()
	rv, gv, bv, av := v.RGBA()
	return ru == rv && bu == bv && gu == gv && au == av
}

func randImage(w, h int) *image.RGBA {
	im := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			r := uint8(rand.Int())
			g := uint8(rand.Int())
			b := uint8(rand.Int())
			im.Set(i, j, color.RGBA{r, g, b, 0xFF})
		}
	}
	return im
}

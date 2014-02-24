package detect

import (
	"image"
	"testing"
)

func TestSuppressOverlap(t *testing.T) {
	cases := []struct {
		MaxInter float64
		MaxNum   int
		In, Out  []Det
	}{
		// Clear margin between windows.
		{
			0, 10,
			[]Det{
				{4, image.Rect(0, 0, 5, 5)},
				{3, image.Rect(10, 0, 15, 5)},
				{2, image.Rect(0, 10, 5, 15)},
				{1, image.Rect(10, 10, 15, 15)},
			},
			[]Det{
				{4, image.Rect(0, 0, 5, 5)},
				{3, image.Rect(10, 0, 15, 5)},
				{2, image.Rect(0, 10, 5, 15)},
				{1, image.Rect(10, 10, 15, 15)},
			},
		},
		// Same and limit to four outputs.
		{
			0, 4,
			[]Det{
				{4, image.Rect(0, 0, 5, 5)},
				{3, image.Rect(10, 0, 15, 5)},
				{2, image.Rect(0, 10, 5, 15)},
				{1, image.Rect(10, 10, 15, 15)},
			},
			[]Det{
				{4, image.Rect(0, 0, 5, 5)},
				{3, image.Rect(10, 0, 15, 5)},
				{2, image.Rect(0, 10, 5, 15)},
				{1, image.Rect(10, 10, 15, 15)},
			},
		},
		// Same and limit to three outputs.
		{
			0, 3,
			[]Det{
				{4, image.Rect(0, 0, 5, 5)},
				{3, image.Rect(10, 0, 15, 5)},
				{2, image.Rect(0, 10, 5, 15)},
				{1, image.Rect(10, 10, 15, 15)},
			},
			[]Det{
				{4, image.Rect(0, 0, 5, 5)},
				{3, image.Rect(10, 0, 15, 5)},
				{2, image.Rect(0, 10, 5, 15)},
			},
		},
		// Touching but not overlapping.
		{
			0, 10,
			[]Det{
				{4, image.Rect(0, 0, 5, 5)},
				{3, image.Rect(5, 0, 10, 5)},
				{2, image.Rect(0, 5, 5, 10)},
				{1, image.Rect(5, 5, 10, 10)},
			},
			[]Det{
				{4, image.Rect(0, 0, 5, 5)},
				{3, image.Rect(5, 0, 10, 5)},
				{2, image.Rect(0, 5, 5, 10)},
				{1, image.Rect(5, 5, 10, 10)},
			},
		},
		// All slightly overlapping.
		{
			0, 10,
			[]Det{
				{4, image.Rect(0, 0, 5, 5)},
				{3, image.Rect(4, 0, 9, 5)},
				{2, image.Rect(0, 4, 5, 9)},
				{1, image.Rect(4, 4, 9, 9)},
			},
			[]Det{
				{4, image.Rect(0, 0, 5, 5)},
			},
		},
		// B and C overlapping A and D. Output A and D.
		{
			0, 10,
			[]Det{
				{4, image.Rect(0, 0, 5, 5)},
				{3, image.Rect(3, 0, 15, 15)},
				{2, image.Rect(0, 3, 15, 15)},
				{1, image.Rect(10, 10, 15, 15)},
			},
			[]Det{
				{4, image.Rect(0, 0, 5, 5)},
				{1, image.Rect(10, 10, 15, 15)},
			},
		},
		// Same, limit to two outputs.
		{
			0, 2,
			[]Det{
				{4, image.Rect(0, 0, 5, 5)},
				{3, image.Rect(3, 0, 15, 15)},
				{2, image.Rect(0, 3, 15, 15)},
				{1, image.Rect(10, 10, 15, 15)},
			},
			[]Det{
				{4, image.Rect(0, 0, 5, 5)},
				{1, image.Rect(10, 10, 15, 15)},
			},
		},
		// Test intersection threshold.
		{
			0.5, 10,
			[]Det{
				{5, image.Rect(0, 0, 5, 5)},
				{4, image.Rect(1, 0, 6, 5)},
				{3, image.Rect(2, 0, 7, 5)},
				{2, image.Rect(3, 0, 8, 5)},
				{1, image.Rect(4, 0, 9, 5)},
			},
			[]Det{
				{5, image.Rect(0, 0, 5, 5)},
				{2, image.Rect(3, 0, 8, 5)},
			},
		},
		{
			0.1, 10,
			[]Det{
				{5, image.Rect(0, 0, 5, 5)},
				{4, image.Rect(1, 0, 6, 5)},
				{3, image.Rect(2, 0, 7, 5)},
				{2, image.Rect(3, 0, 8, 5)},
				{1, image.Rect(4, 0, 9, 5)},
			},
			[]Det{
				{5, image.Rect(0, 0, 5, 5)},
			},
		},
		{
			0.3, 10,
			[]Det{
				{5, image.Rect(0, 0, 5, 5)},
				{4, image.Rect(1, 0, 6, 5)},
				{3, image.Rect(2, 0, 7, 5)},
				{2, image.Rect(3, 0, 8, 5)},
				{1, image.Rect(4, 0, 9, 5)},
			},
			[]Det{
				{5, image.Rect(0, 0, 5, 5)},
				{1, image.Rect(4, 0, 9, 5)},
			},
		},
		{
			0.6 + 0.005, 10,
			[]Det{
				{5, image.Rect(0, 0, 5, 5)},
				{4, image.Rect(1, 0, 6, 5)},
				{3, image.Rect(2, 0, 7, 5)},
				{2, image.Rect(3, 0, 8, 5)},
				{1, image.Rect(4, 0, 9, 5)},
			},
			[]Det{
				{5, image.Rect(0, 0, 5, 5)},
				{3, image.Rect(2, 0, 7, 5)},
				{1, image.Rect(4, 0, 9, 5)},
			},
		},
		// Same test but vertical.
		{
			0.6 + 0.005, 10,
			[]Det{
				{5, image.Rect(0, 0, 5, 5)},
				{4, image.Rect(0, 1, 5, 6)},
				{3, image.Rect(0, 2, 5, 7)},
				{2, image.Rect(0, 3, 5, 8)},
				{1, image.Rect(0, 4, 5, 9)},
			},
			[]Det{
				{5, image.Rect(0, 0, 5, 5)},
				{3, image.Rect(0, 2, 5, 7)},
				{1, image.Rect(0, 4, 5, 9)},
			},
		},
		// Same test but diagonal.
		{
			0.5 * 0.5, 10,
			[]Det{
				{5, image.Rect(0, 0, 5, 5)},
				{4, image.Rect(1, 1, 6, 6)},
				{3, image.Rect(2, 2, 7, 7)},
				{2, image.Rect(3, 3, 8, 8)},
				{1, image.Rect(4, 4, 9, 9)},
			},
			[]Det{
				{5, image.Rect(0, 0, 5, 5)},
				{2, image.Rect(3, 3, 8, 8)},
			},
		},
	}

	for _, x := range cases {
		out := SuppressOverlap(x.In, x.MaxNum, x.MaxInter)
		if len(out) != len(x.Out) {
			t.Error("different length")
			t.Log(x)
			t.Log(out)
			continue
		}
		for i := range x.Out {
			if !x.Out[i].Pos.Eq(out[i].Pos) {
				t.Errorf("differ at index %d", i)
				t.Log(x)
				t.Logf("want: %v", out)
				t.Logf("got:  %v", x.Out)
				break
			}
		}
	}
}

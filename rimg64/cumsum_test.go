package rimg64

import (
	"image"
	"math"
	"testing"
)

func TestCumSum(t *testing.T) {
	f := FromRows([][]float64{
		[]float64{1, 2, 3, 4},
		[]float64{-1, 2, -3, 4},
		[]float64{0.5, 0.5, 0.5, 0.5},
	})
	want := FromRows([][]float64{
		[]float64{1, 3, 6, 10},
		[]float64{0, 4, 4, 12},
		[]float64{0.5, 5, 5.5, 14},
	})
	testEq(t, want, (*Image)(CumSum(f)), 1e-16)
}

func TestCumSumRect(t *testing.T) {
	const eps = 1e-16
	f := FromRows([][]float64{
		[]float64{1, 2, 3, 4},
		[]float64{-1, 2, -3, 4},
		[]float64{0.5, 0.5, 0.5, 0.5},
	})
	table := CumSum(f)
	tests := []struct {
		Rect image.Rectangle
		Want float64
	}{
		{image.Rectangle{image.Pt(0, 0), image.Pt(1, 1)}, 1},
		{image.Rectangle{image.Pt(0, 1), image.Pt(1, 2)}, -1},
		{image.Rectangle{image.Pt(0, 0), image.Pt(4, 3)}, 14},
		{image.Rectangle{image.Pt(3, 2), image.Pt(4, 3)}, 0.5},
		{image.Rectangle{image.Pt(1, 1), image.Pt(3, 3)}, 2 - 3 + 0.5 + 0.5},
		{image.Rectangle{image.Pt(0, 0), image.Pt(3, 1)}, 1 + 2 + 3},
		{image.Rectangle{image.Pt(1, 0), image.Pt(4, 1)}, 2 + 3 + 4},
		{image.Rectangle{image.Pt(0, 2), image.Pt(3, 3)}, 0.5 + 0.5 + 0.5},
		{image.Rectangle{image.Pt(0, 0), image.Pt(0, 0)}, 0},
		{image.Rectangle{image.Pt(3, 2), image.Pt(3, 2)}, 0},
		{image.Rectangle{image.Pt(0, 2), image.Pt(3, 2)}, 0},
		{image.Rectangle{image.Pt(1, 0), image.Pt(1, 4)}, 0},
	}
	for _, test := range tests {
		got := table.Rect(test.Rect)
		if math.Abs(got-test.Want) > eps {
			t.Errorf("error: rect %v: want %g, got %g", test.Rect, test.Want, got)
		}
	}
}

func testEq(t *testing.T, want, got *Image, eps float64) bool {
	if !want.Size().Eq(got.Size()) {
		t.Errorf("different size: want %v, got %v", want.Size(), got.Size())
		return false
	}

	eq := true
	for i := 0; i < want.Width; i++ {
		for j := 0; j < want.Height; j++ {
			x := want.At(i, j)
			y := got.At(i, j)
			if math.Abs(x-y) > eps {
				t.Errorf("different: at %d, %d: want %g, got %g", i, j, x, y)
				eq = false
			}
		}
	}
	return eq
}

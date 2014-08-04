package rimg64

import (
	"math"
	"testing"
)

func TestCumSum(t *testing.T) {
	f := FromRows([][]float64 {
		[]float64{1, 2, 3, 4},
		[]float64{-1, 2, -3, 4},
		[]float64{0.5, 0.5, 0.5, 0.5},
	})
	want := FromRows([][]float64 {
		[]float64{1, 3, 6, 10},
		[]float64{0, 4, 4, 12},
		[]float64{0.5, 5, 5.5, 14},
	})
	testEq(t, want, CumSum(f), 1e-16)
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
			if math.Abs(x - y) > eps {
				t.Errorf("different: at %d, %d: want %g, got %g", i, j, x, y)
				eq = false
			}
		}
	}
	return eq
}

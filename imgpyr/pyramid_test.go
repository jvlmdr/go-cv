package imgpyr

import (
	"image"
	"math"
	"testing"
)

func TestGeoSeq_Reverse(t *testing.T) {
	const (
		n   = 4
		eps = 1e-12
	)
	seq := GeoSeq{1, 2, n}
	rev := seq.Reverse()
	for i := 0; i < seq.Len; i++ {
		got := rev.At(i)
		want := seq.At(n - 1 - i)
		if math.Abs(want-got) > eps {
			t.Errorf("wrong value at %d (want %g, got %g)", i, want, got)
		}
	}
}

func TestSequence(t *testing.T) {
	var (
		start float64 = 2
		step  float64 = 2
		max   float64 = 35
	)
	seq := Sequence(start, step, max)
	last := seq.At(seq.Len - 1)
	if last > max {
		t.Errorf("expected last <= max (%g <= %g)", last, max)
	}
	if step*last <= max {
		t.Errorf("expected step * last > max (%g < %g)", step*last, max)
	}
}

func TestLogRange(t *testing.T) {
	var (
		first float64 = 2
		last  float64 = 32
		n     int     = 10
		eps   float64 = 1e-12
	)
	seq := LogRange(first, last, n)
	got := seq.At(seq.Len - 1)
	if math.Abs(got-last) > eps {
		t.Errorf("wrong last element (got %g, want %g)", got, last)
	}
}

func TestScales(t *testing.T) {
	cases := []struct {
		Image, Tmpl image.Point
		Step        float64
		Len         int
	}{
		{image.Pt(400, 300), image.Pt(40, 30), math.Sqrt(10), 3},
		{image.Pt(400, 300), image.Pt(40, 40), math.Sqrt(10), 2},
		{image.Pt(400, 300), image.Pt(0, 30), math.Sqrt(10), 3},
		{image.Pt(512, 256), image.Pt(1, 1), 2, 9},
	}

	for _, e := range cases {
		scales := Scales(e.Image, e.Tmpl, e.Step)
		if scales.Len != e.Len {
			t.Errorf(
				"image %v, min %v, step %g: want %d, got %d",
				e.Image, e.Tmpl, e.Step, e.Len, scales.Len,
			)
		}
	}
}

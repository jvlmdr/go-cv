package detect_test

import (
	"testing"

	"github.com/jvlmdr/go-cv/detect"
)

func TestMissRateAtFPPI(t *testing.T) {
	valset := &detect.ValSet{
		Dets: []detect.ValScore{
			{0, true},
			{-1, true},
			// < 1 FP: 3 misses
			{-3, false},
			{-4, true},
			{-5, true},
			// < 2 FP: 1 miss
			{-5, false},
			// < 3 FP: 1 miss
			{-6, false},
			{-7, true},
			// < 4 FP: 0 misses
			{-8, false},
			{-9, false},
		},
		Misses: 7,
		Images: 13,
	}
	const numPos float64 = 5 + 7
	cases := []struct {
		FP       int
		MissRate float64
	}{
		{0, (3 + 7) / numPos},
		{1, (1 + 7) / numPos},
		{2, (1 + 7) / numPos},
		{3, (0 + 7) / numPos},
		{4, (0 + 7) / numPos},
		{13, (0 + 7) / numPos},
		{100, (0 + 7) / numPos},
	}
	for _, s := range cases {
		fppi := (float64(s.FP) + 0.5) / float64(valset.Images)
		rate := detect.MissRateAtFPPI(valset, fppi)
		if rate != s.MissRate {
			t.Errorf("false positives %d: want %.3g, got %.3g", s.FP, s.MissRate, rate)
		}
	}
}

package slide_test

import (
	"image"
	"testing"

	"github.com/jvlmdr/go-cv/slide"
)

func TestCorrMultiStrideNaive_vsDecimate(t *testing.T) {
	const eps = 1e-9
	cases := []struct {
		ImSize   image.Point
		TmplSize image.Point
		C        int
		K        int
	}{
		{ImSize: image.Pt(8, 10), TmplSize: image.Pt(3, 2), C: 5, K: 5},
		{ImSize: image.Pt(100, 1), TmplSize: image.Pt(1, 1), C: 5, K: 5},
		{ImSize: image.Pt(1, 100), TmplSize: image.Pt(1, 1), C: 5, K: 5},
		{ImSize: image.Pt(43, 64), TmplSize: image.Pt(4, 5), C: 5, K: 3},
		{ImSize: image.Pt(43, 64), TmplSize: image.Pt(5, 4), C: 5, K: 3},
		{ImSize: image.Pt(64, 43), TmplSize: image.Pt(4, 5), C: 5, K: 3},
		{ImSize: image.Pt(64, 43), TmplSize: image.Pt(5, 4), C: 5, K: 3},
		{ImSize: image.Pt(63, 127), TmplSize: image.Pt(3, 2), C: 5, K: 32},
		{ImSize: image.Pt(63, 127), TmplSize: image.Pt(2, 3), C: 5, K: 32},
		{ImSize: image.Pt(63, 127), TmplSize: image.Pt(3, 2), C: 5, K: 31},
		{ImSize: image.Pt(63, 127), TmplSize: image.Pt(2, 3), C: 5, K: 31},
		{ImSize: image.Pt(63, 127), TmplSize: image.Pt(2, 3), C: 5, K: 10000},
	}

	for _, q := range cases {
		f := randMulti(q.ImSize.X, q.ImSize.Y, q.C)
		g := randMulti(q.TmplSize.X, q.TmplSize.Y, q.C)
		h, err := slide.CorrMulti(f, g)
		if err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.ImSize, q.TmplSize, q.K, err)
			continue
		}
		want := slide.Decimate(h, q.K)
		got, err := slide.CorrMultiStrideNaive(f, g, q.K)
		if err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.ImSize, q.TmplSize, q.K, err)
			continue
		}
		if err := errIfNotEqImage(want, got, eps); err != nil {
			t.Errorf("im %v, tmpl %v, stride %d: %v", q.ImSize, q.TmplSize, q.K, err)
			continue
		}
	}
}

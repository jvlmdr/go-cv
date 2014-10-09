package convfeat_test

import (
	"math"
	"testing"

	"github.com/jvlmdr/go-cv/convfeat"
	"github.com/jvlmdr/go-cv/rimg64"
)

func TestAdjChanNorm_Apply(t *testing.T) {
	const eps = 1e-9
	const (
		k = 2
		n = 5
		a = 1e-4
		b = 0.75
	)
	phi := &convfeat.AdjChanNorm{K: k, Num: n, Alpha: a, Beta: b}
	elems := []float64{-1, 2, 3, 2, 1, -1, -3}
	cases := []struct {
		In, Out []float64
	}{{
		In: elems,
		Out: []float64{
			elems[0] / math.Pow(k+a*sumsqr(elems[:3]), b),
			elems[1] / math.Pow(k+a*sumsqr(elems[:4]), b),
			elems[2] / math.Pow(k+a*sumsqr(elems[:5]), b),
			elems[3] / math.Pow(k+a*sumsqr(elems[1:6]), b),
			elems[4] / math.Pow(k+a*sumsqr(elems[2:]), b),
			elems[5] / math.Pow(k+a*sumsqr(elems[3:]), b),
			elems[6] / math.Pow(k+a*sumsqr(elems[4:]), b),
		},
	}, {
		In: elems[:6],
		Out: []float64{
			elems[0] / math.Pow(k+a*sumsqr(elems[:6][:3]), b),
			elems[1] / math.Pow(k+a*sumsqr(elems[:6][:4]), b),
			elems[2] / math.Pow(k+a*sumsqr(elems[:6][:5]), b),
			elems[3] / math.Pow(k+a*sumsqr(elems[:6][1:]), b),
			elems[4] / math.Pow(k+a*sumsqr(elems[:6][2:]), b),
			elems[5] / math.Pow(k+a*sumsqr(elems[:6][3:]), b),
		},
	}, {
		In: elems[:5],
		Out: []float64{
			elems[0] / math.Pow(k+a*sumsqr(elems[:5][:3]), b),
			elems[1] / math.Pow(k+a*sumsqr(elems[:5][:4]), b),
			elems[2] / math.Pow(k+a*sumsqr(elems[:5]), b),
			elems[3] / math.Pow(k+a*sumsqr(elems[:5][1:]), b),
			elems[4] / math.Pow(k+a*sumsqr(elems[:5][2:]), b),
		},
	}, {
		In: elems[:3],
		Out: []float64{
			elems[0] / math.Pow(k+a*sumsqr(elems[:3]), b),
			elems[1] / math.Pow(k+a*sumsqr(elems[:3]), b),
			elems[2] / math.Pow(k+a*sumsqr(elems[:3]), b),
		},
	}, {
		In: elems[:2],
		Out: []float64{
			elems[0] / math.Pow(k+a*sumsqr(elems[:2]), b),
			elems[1] / math.Pow(k+a*sumsqr(elems[:2]), b),
		},
	}, {
		In: elems[:1],
		Out: []float64{
			elems[0] / math.Pow(k+a*sumsqr(elems[:1]), b),
		},
	}}

	for _, test := range cases {
		f := rimg64.NewMulti(1, 1, len(test.In))
		f.SetPixel(0, 0, test.In)
		y, err := phi.Apply(f)
		if err != nil {
			t.Fatal(err)
		}
		for i := range test.Out {
			want, got := test.Out[i], y.At(0, 0, i)
			if math.Abs(want-got) > eps {
				t.Errorf("with %d channels: different at %d: want %g, got %g", len(test.In), i, want, got)
			}
		}
	}
}

func sumsqr(x []float64) float64 {
	var t float64
	for _, e := range x {
		t += e * e
	}
	return t
}

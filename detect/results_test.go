package detect_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/jvlmdr/go-cv/detect"
	"github.com/jvlmdr/go-ml/ml"
)

func TestValSet_Enum(t *testing.T) {
	cases := []struct {
		In  *detect.ValSet
		Out ml.PerfPath
	}{
		// Empty.
		{
			&detect.ValSet{
				Dets:   []detect.ValScore{},
				Misses: 0,
			},
			[]ml.Perf{
				{TP: 0, FP: 0, FN: 0},
			},
		},
		// Empty with misses.
		{
			&detect.ValSet{
				Dets:   []detect.ValScore{},
				Misses: 8,
			},
			[]ml.Perf{
				{TP: 0, FP: 0, FN: 8},
			},
		},
		// One true detection.
		{
			&detect.ValSet{
				Dets:   []detect.ValScore{{True: true}},
				Misses: 8,
			},
			[]ml.Perf{
				{TP: 0, FP: 0, FN: 9},
				{TP: 1, FP: 0, FN: 8},
			},
		},
		// One false detection.
		{
			&detect.ValSet{
				Dets:   []detect.ValScore{{True: false}},
				Misses: 8,
			},
			[]ml.Perf{
				{TP: 0, FP: 0, FN: 8},
				{TP: 0, FP: 1, FN: 8},
			},
		},
		// Whole string of detections.
		{
			&detect.ValSet{
				Dets: []detect.ValScore{
					{True: true},
					{True: true},
					{True: false},
					{True: true},
					{True: false},
					{True: false},
					{True: false},
					{True: true},
					{True: false},
					{True: false},
					{True: false},
					{True: false},
					{True: false},
					{True: false},
				},
				Misses: 8,
			},
			[]ml.Perf{
				{TP: 0, FP: 0, FN: 12},
				{TP: 1, FP: 0, FN: 11},
				{TP: 2, FP: 0, FN: 10},
				{TP: 2, FP: 1, FN: 10},
				{TP: 3, FP: 1, FN: 9},
				{TP: 3, FP: 2, FN: 9},
				{TP: 3, FP: 3, FN: 9},
				{TP: 3, FP: 4, FN: 9},
				{TP: 4, FP: 4, FN: 8},
				{TP: 4, FP: 5, FN: 8},
				{TP: 4, FP: 6, FN: 8},
				{TP: 4, FP: 7, FN: 8},
				{TP: 4, FP: 8, FN: 8},
				{TP: 4, FP: 9, FN: 8},
				{TP: 4, FP: 10, FN: 8},
			},
		},
	}

	for _, x := range cases {
		out := x.In.Enum()
		if !reflect.DeepEqual(out, x.Out) {
			s := fmt.Sprint(
				"detections:\n\t", x.In.Dets, "\n",
				"misses: ", x.In.Misses, "\n",
				"want:\n\t", x.Out, "\n",
				"got:\n\t", out,
			)
			t.Error("different results\n" + s)
		}
	}
}

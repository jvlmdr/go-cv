package detect

import (
	"github.com/jackvalmadre/go-ml"

	"fmt"
	"reflect"
	"testing"
)

func TestResultSet_Enum(t *testing.T) {
	cases := []struct {
		In  *ResultSet
		Out []ml.Result
	}{
		// Empty.
		{
			&ResultSet{
				Dets:   []ValDet{},
				Misses: 0,
			},
			[]ml.Result{
				{TP: 0, FP: 0, FN: 0},
			},
		},
		// Empty with misses.
		{
			&ResultSet{
				Dets:   []ValDet{},
				Misses: 8,
			},
			[]ml.Result{
				{TP: 0, FP: 0, FN: 8},
			},
		},
		// One true detection.
		{
			&ResultSet{
				Dets:   []ValDet{{True: true}},
				Misses: 8,
			},
			[]ml.Result{
				{TP: 0, FP: 0, FN: 9},
				{TP: 1, FP: 0, FN: 8},
			},
		},
		// One false detection.
		{
			&ResultSet{
				Dets:   []ValDet{{True: false}},
				Misses: 8,
			},
			[]ml.Result{
				{TP: 0, FP: 0, FN: 8},
				{TP: 0, FP: 1, FN: 8},
			},
		},
		// Whole string of detections.
		{
			&ResultSet{
				Dets: []ValDet{
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
			[]ml.Result{
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

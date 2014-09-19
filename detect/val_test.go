package detect_test

import (
	"fmt"
	"image"
	"reflect"
	"testing"

	"github.com/jvlmdr/go-cv/detect"
)

func TestMatch(t *testing.T) {
	cases := []struct {
		Dets     []detect.Det
		Refs     []image.Rectangle
		MinInter float64
		Match    map[int]int
	}{
		// Empty.
		{
			[]detect.Det{},
			[]image.Rectangle{},
			0.5,
			map[int]int{},
		},
		// No detections.
		{
			[]detect.Det{},
			[]image.Rectangle{
				image.Rect(10, 10, 110, 110),
				image.Rect(100, 0, 200, 100),
			},
			0.5,
			map[int]int{},
		},
		// No references.
		{
			[]detect.Det{
				{10, image.Rect(0, 0, 100, 100)},
				{9, image.Rect(110, 10, 210, 110)},
			},
			[]image.Rectangle{},
			0.5,
			map[int]int{},
		},

		// One detection, two references, one match.
		// A: 40 * 50 = 2000
		// B: 40 * 60 = 2400
		// A cap B: (40-10) * (70-20)) = 30 * 50 = 1500
		// A cup B: 2000 + 2400 - 1500 = 2900
		// 1500 / 2900 > 0.5
		{
			[]detect.Det{
				{10, image.Rect(10, 20, 50, 70)},
			},
			[]image.Rectangle{
				image.Rect(0, 10, 40, 70),
				image.Rect(90, 10, 120, 40),
			},
			0.5,
			map[int]int{0: 0},
		},
		// Different order of references.
		{
			[]detect.Det{
				{10, image.Rect(10, 20, 50, 70)},
			},
			[]image.Rectangle{
				image.Rect(90, 10, 120, 40),
				image.Rect(0, 10, 40, 70),
			},
			0.5,
			map[int]int{0: 1},
		},
		// One detection, two references, no matches.
		{
			[]detect.Det{
				{10, image.Rect(10, 20, 50, 80)},
			},
			[]image.Rectangle{
				image.Rect(0, 90, 40, 160),
				image.Rect(90, 10, 120, 40),
			},
			0.5,
			map[int]int{},
		},

		// Check intersection threshold.
		// (100-33) / 133 > 0.5
		{
			[]detect.Det{
				{10, image.Rect(0, 0, 100, 100)},
			},
			[]image.Rectangle{
				image.Rect(33, 0, 133, 100),
			},
			0.5,
			map[int]int{0: 0},
		},
		// (100-34) / 134 < 0.5
		{
			[]detect.Det{
				{10, image.Rect(0, 0, 100, 100)},
			},
			[]image.Rectangle{
				image.Rect(34, 0, 134, 100),
			},
			0.5,
			map[int]int{},
		},
		// (100-50) / 150 = 1/3
		{
			[]detect.Det{
				{10, image.Rect(0, 0, 100, 100)},
			},
			[]image.Rectangle{
				image.Rect(50, 0, 150, 100),
			},
			1.0/3.0 - 0.005,
			map[int]int{0: 0},
		},
		// (100-50) / 150 = 1/3
		{
			[]detect.Det{
				{10, image.Rect(0, 0, 100, 100)},
			},
			[]image.Rectangle{
				image.Rect(50, 0, 150, 100),
			},
			1.0/3.0 + 0.005,
			map[int]int{},
		},

		// Match first to first and second to second.
		{
			[]detect.Det{
				{10, image.Rect(0, 0, 100, 100)},
				{9, image.Rect(110, 10, 210, 110)},
			},
			[]image.Rectangle{
				image.Rect(10, 10, 110, 110),
				image.Rect(100, 0, 200, 100),
			},
			0.5,
			map[int]int{0: 0, 1: 1},
		},
		// Match first to second and second to first.
		{
			[]detect.Det{
				{10, image.Rect(0, 0, 100, 100)},
				{9, image.Rect(110, 10, 210, 110)},
			},
			[]image.Rectangle{
				image.Rect(100, 0, 200, 100),
				image.Rect(10, 10, 110, 110),
			},
			0.5,
			map[int]int{0: 1, 1: 0},
		},
		// Match first to third even though first is OK.
		{
			[]detect.Det{
				{10, image.Rect(0, 0, 100, 100)},
				{9, image.Rect(110, 10, 210, 110)},
			},
			[]image.Rectangle{
				image.Rect(10, 10, 110, 110),
				image.Rect(100, 0, 200, 100),
				image.Rect(-5, -5, 95, 95),
			},
			0.5,
			map[int]int{0: 2, 1: 1},
		},

		// Let first detection take reference
		// even though it's better for the second.
		{
			[]detect.Det{
				{10, image.Rect(0, 0, 100, 100)},
				{9, image.Rect(5, 5, 105, 105)},
			},
			[]image.Rectangle{
				image.Rect(10, 10, 110, 110),
			},
			0.5,
			map[int]int{0: 0},
		},
		// Let first detection take reference
		// which is better for the second.
		// Provide a reference for the second detection.
		{
			[]detect.Det{
				{10, image.Rect(0, 0, 100, 100)},
				{9, image.Rect(5, 5, 105, 105)},
			},
			[]image.Rectangle{
				image.Rect(15, 15, 115, 115),
				image.Rect(10, 10, 110, 110),
				image.Rect(5, 5, 105, 105),
			},
			0.5,
			map[int]int{0: 2, 1: 1},
		},
	}

	for _, x := range cases {
		match := detect.Match(x.Dets, x.Refs, x.MinInter)
		if !reflect.DeepEqual(match, x.Match) {
			s := fmt.Sprint(
				"detections:\n\t", x.Dets, "\n",
				"references:\n\t", x.Refs, "\n",
				"min. inter: ", x.MinInter, "\n",
				"want:\n\t", x.Match, "\n",
				"got:\n\t", match,
			)
			t.Error("different matches\n" + s)
		}
	}
}

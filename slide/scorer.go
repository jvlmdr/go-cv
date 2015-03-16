package slide

import (
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
)

// Scorer is a method for scoring images of a fixed size.
type Scorer interface {
	Size() image.Point
	Score(*rimg64.Multi) (float64, error)
}

// Slider is a Scorer that has an efficient method
// for evaluating itself in sliding window fashion.
type Slider interface {
	Scorer
	Slide(*rimg64.Multi) (*rimg64.Image, error)
}

// Score computes the score of every window.
// If scorer is a Slider, then its Slide() function is called.
func Score(im *rimg64.Multi, scorer Scorer) (*rimg64.Image, error) {
	var err error
	var resp *rimg64.Image
	if slider, ok := scorer.(Slider); ok {
		resp, err = slider.Slide(im)
		if err != nil {
			return nil, err
		}
	} else {
		resp, err = EvalFunc(im, scorer.Size(), scorer.Score)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}

// It could be useful to move beyond simply Scorer and Slider.
// For example, to evaluate multiple detectors on the same image
// the Fourier transform and integral images could be re-used.

//	type SliderList interface {
//		Len() int
//		At(int) Scorer
//		Slide(*rimg64.Multi) ([]*rimg64.Image, error)
//	}

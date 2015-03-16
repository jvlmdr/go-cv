package slide

import (
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
)

// ScoreFunc scores a multi-channel image.
type ScoreFunc func(*rimg64.Multi) (float64, error)

// EvalFunc evaluates a function on every window in an image.
// If the input image is M x N and the window size is m x n,
// then the output is (M-m+1) x (N-n+1).
// If the window size is larger than the image size in either dimension,
// a nil image is returned with no error.
func EvalFunc(im *rimg64.Multi, size image.Point, f ScoreFunc) (*rimg64.Image, error) {
	if im.Width < size.X || im.Height < size.Y {
		return nil, nil
	}
	r := rimg64.New(im.Width-size.X+1, im.Height-size.Y+1)
	x := rimg64.NewMulti(size.X, size.Y, im.Channels)
	for i := 0; i < r.Width; i++ {
		for j := 0; j < r.Height; j++ {
			// Copy window into x.
			for u := 0; u < size.X; u++ {
				for v := 0; v < size.Y; v++ {
					for p := 0; p < im.Channels; p++ {
						x.Set(u, v, p, im.At(i+u, j+v, p))
					}
				}
			}
			y, err := f(x)
			if err != nil {
				return nil, err
			}
			r.Set(i, j, y)
		}
	}
	return r, nil
}

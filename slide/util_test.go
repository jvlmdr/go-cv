package slide_test

import (
	"math/rand"

	"github.com/jvlmdr/go-cv/rimg64"
)

func randImage(width, height int) *rimg64.Image {
	f := rimg64.New(width, height)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			f.Set(i, j, rand.NormFloat64())
		}
	}
	return f
}

func randMulti(width, height, channels int) *rimg64.Multi {
	f := rimg64.NewMulti(width, height, channels)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			for k := 0; k < channels; k++ {
				f.Set(i, j, k, rand.NormFloat64())
			}
		}
	}
	return f
}

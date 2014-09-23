package feat

import (
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
)

// NewImageFunc creates a transform which calls the given function.
func NewImageFunc(f func(image.Image) (*rimg64.Multi, error), rate int) Image {
	return &imageFunc{f, rate}
}

type imageFunc struct {
	apply func(image.Image) (*rimg64.Multi, error)
	rate  int
}

func (t *imageFunc) Apply(im image.Image) (*rimg64.Multi, error) {
	return t.apply(im)
}

func (t *imageFunc) Rate() int {
	return t.rate
}

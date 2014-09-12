package imgpyr

import (
	"image"

	"github.com/nfnt/resize"
)

// Generator describes an image pyramid which is generated level by level.
//
// The scales must be monotonically decreasing.
type Generator struct {
	Image  image.Image
	Scales []float64
	Interp resize.InterpolationFunction
}

// Level describes one level of a pyramid Generator.
type Level struct {
	Image image.Image
	Index int
}

func NewGenerator(im image.Image, scales []float64, interp resize.InterpolationFunction) *Generator {
	if len(scales) == 0 {
		return nil
	}
	return &Generator{im, scales, interp}
}

func (pyr *Generator) First() *Level {
	if len(pyr.Scales) == 0 {
		return nil
	}
	size := scaleSize(pyr.Image.Bounds().Size(), pyr.Scales[0])
	im := resizeIfNec(size, pyr.Image, pyr.Interp)
	return &Level{im, 0}
}

func (pyr *Generator) Next(curr *Level) *Level {
	index := curr.Index + 1
	if index >= len(pyr.Scales) {
		// This was the last scale.
		return nil
	}
	var src image.Image
	if pyr.Scales[curr.Index] >= 1 {
		// The image at the current level was larger than the original.
		// Rescale the original.
		src = pyr.Image
	} else {
		// The image at the current level was smaller than the original.
		// Rescale the current level.
		src = curr.Image
	}
	size := scaleSize(pyr.Image.Bounds().Size(), pyr.Scales[index])
	im := resizeIfNec(size, src, pyr.Interp)
	return &Level{im, index}
}

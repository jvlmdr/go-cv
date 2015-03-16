package detect

import "github.com/jvlmdr/go-cv/slide"

// FeatTmpl is an affine template.
type FeatTmpl struct {
	// Assigns a score to feature images of a fixed size.
	Scorer *slide.AffineScorer
	// The size of the image from which the features were computed,
	// and the position of the bounding box within it.
	PixelShape PadRect
}

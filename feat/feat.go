package feat

import (
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
)

// Transform defines a common interface for feature transforms.
//
// Feature transforms are assumed to have a positive integer downsample rate.
// If the input image f(x, y) has size (m, n) with domain
// 	x = 0, ..., m - 1
// 	y = 0, ..., n - 1
// and produces a feature image g(u, v) of size (p, q) with domain
// 	u = 0, ..., p - 1
// 	v = 0, ..., q - 1
// then calling Apply() on any sub-image
// 	x = rate*left, ..., m - 1 - rate*right
// 	y = rate*top, ..., n - 1 - rate*bottom
// must produce the feature image
// 	u = left, ..., p - 1 - right
// 	v = top, ..., q - 1 - bottom
// where (left, right, top, bottom) are non-negative integers describing an inset on each side.
type Transform interface {
	// Function to compute transform on image.
	Apply(im image.Image) *rimg64.Multi
	// Integer downsample rate.
	Rate() int
}

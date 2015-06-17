package feat

import (
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
)

// Image defines a transform of an image.Image.
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
type Image interface {
	// Function to compute transform on image.
	Apply(image.Image) (*rimg64.Multi, error)
	// Integer downsample rate.
	Rate() int
	// The size of the feature image computed from an input image.
	Size(im image.Point) (feat image.Point)
	// Minimum image size to achieve feature image size.
	MinInputSize(feat image.Point) (im image.Point)
	// The number of channels.
	Channels() int
}

// Real defines a transform of a real-valued image.
type Real interface {
	// Function to compute transform on image.
	Apply(*rimg64.Multi) (*rimg64.Multi, error)
	// Integer downsample rate.
	Rate() int
	// The size of the feature image computed from an input image.
	Size(im image.Point) (feat image.Point)
	// Minimum image size to achieve feature image size.
	MinInputSize(feat image.Point) (im image.Point)
	// The number of channels.
	Channels() int
}

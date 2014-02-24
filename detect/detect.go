package detect

import "image"

// Detection in image.
type Det struct {
	Score float64
	Rect  image.Rectangle
}

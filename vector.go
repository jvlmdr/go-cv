package cv

import (
	"github.com/jackvalmadre/go-vec"
)

// Describes a vectorized real image.
type RealImageAsVector struct {
	Image RealImage
}

func (vec RealImageAsVector) Size() int {
	image := &vec.Image
	return image.Width * image.Height
}

func (vec RealImageAsVector) At(i int) float64 {
	return vec.Image.Pixels[i]
}

func (vec RealImageAsVector) Set(i int, v float64) {
	vec.Image.Pixels[i] = v
}

func (vec RealImageAsVector) Type() vec.Type {
	image := &vec.Image
	return RealImageAsVectorType{image.Width, image.Height}
}

type RealImageAsVectorType struct {
	Width  int
	Height int
}

func (t RealImageAsVectorType) Size() int {
	return t.Width * t.Height
}

func (t RealImageAsVectorType) New() vec.MutableTyped {
	image := NewRealImage(t.Width, t.Height)
	return RealImageAsVector{image}
}

// Describes a vectorized image of real vectors.
type RealVectorImageAsVector struct {
	Image RealVectorImage
}

func (vec RealVectorImageAsVector) Size() int {
	image := &vec.Image
	return image.Width * image.Height * image.Channels
}

func (vec RealVectorImageAsVector) At(i int) float64 {
	return vec.Image.Pixels[i]
}

func (vec RealVectorImageAsVector) Set(i int, v float64) {
	vec.Image.Pixels[i] = v
}

func (vec RealVectorImageAsVector) Type() vec.Type {
	image := &vec.Image
	return RealVectorImageAsVectorType{image.Width, image.Height, image.Channels}
}

// Describes the type of such a vector.
type RealVectorImageAsVectorType struct {
	Width    int
	Height   int
	Channels int
}

func (t RealVectorImageAsVectorType) Size() int {
	return t.Width * t.Height * t.Channels
}

func (t RealVectorImageAsVectorType) New() vec.MutableTyped {
	image := NewRealVectorImage(t.Width, t.Height, t.Channels)
	return RealVectorImageAsVector{image}
}

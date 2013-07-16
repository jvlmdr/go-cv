package cv

import "github.com/jackvalmadre/lin-go/vec"

// Describes a vectorized real image.
type RealImageAsVector struct {
	Image RealImage
}

func (x RealImageAsVector) Size() int {
	return x.Image.Width * x.Image.Height
}

func (x RealImageAsVector) At(i int) float64 {
	return x.Image.Pixels[i]
}

func (x RealImageAsVector) Set(i int, v float64) {
	x.Image.Pixels[i] = v
}

func (x RealImageAsVector) Type() vec.Type {
	return RealImageAsVectorType{x.Image.Width, x.Image.Height}
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

func (x RealVectorImageAsVector) Size() int {
	return x.Image.Width * x.Image.Height * x.Image.Channels
}

func (x RealVectorImageAsVector) At(i int) float64 {
	return x.Image.Pixels[i]
}

func (x RealVectorImageAsVector) Set(i int, v float64) {
	x.Image.Pixels[i] = v
}

func (x RealVectorImageAsVector) Type() vec.Type {
	return RealVectorImageAsVectorType{x.Image.Width, x.Image.Height, x.Image.Channels}
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

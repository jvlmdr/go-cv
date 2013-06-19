package cv

import (
	"fmt"
	"github.com/jackvalmadre/vector"
	"github.com/skelterjohn/go.matrix"
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

// Describes a vectorized image of real vectors.
type RealVectorImageAsColumnVector struct {
	Image RealVectorImage
}

func (x RealVectorImageAsColumnVector) Nil() bool {
	return x.Image.Empty()
}

func (x RealVectorImageAsColumnVector) Rows() int {
	return x.NumElements()
}

func (x RealVectorImageAsColumnVector) Cols() int {
	return 1
}

func (x RealVectorImageAsColumnVector) NumElements() int {
	f := x.Image
	return f.Width * f.Height * f.Channels
}

func (x RealVectorImageAsColumnVector) GetSize() (int, int) {
	return x.NumElements(), 1
}

func (x RealVectorImageAsColumnVector) Get(i, j int) float64 {
	if j != 0 {
		panic("Out of range")
	}
	return x.Image.Pixels[i]
}

func (x RealVectorImageAsColumnVector) Plus(y matrix.MatrixRO) (matrix.Matrix, error) {
	return x.DenseMatrix().Plus(y)
}

func (x RealVectorImageAsColumnVector) Minus(y matrix.MatrixRO) (matrix.Matrix, error) {
	return x.DenseMatrix().Minus(y)
}

func (x RealVectorImageAsColumnVector) Times(y matrix.MatrixRO) (matrix.Matrix, error) {
	return x.DenseMatrix().Times(y)
}

func (x RealVectorImageAsColumnVector) Det() float64 {
	panic("Not a square matrix")
}

func (x RealVectorImageAsColumnVector) Trace() float64 {
	panic("Not a square matrix")
}

func (x RealVectorImageAsColumnVector) String() string {
	n := x.NumElements()
	return fmt.Sprintf("[%dx1 matrix]", n)
}

func (x RealVectorImageAsColumnVector) DenseMatrix() *matrix.DenseMatrix {
	n := x.NumElements()
	y := matrix.Zeros(n, 1)
	for i := 0; i < n; i++ {
		y.Set(i, 0, x.Get(i, 0))
	}
	return y
}

func (x RealVectorImageAsColumnVector) SparseMatrix() *matrix.SparseMatrix {
	panic("Unlikely (and unimplemented)")
}

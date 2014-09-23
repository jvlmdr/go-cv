package feat

import (
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
)

func init() {
	RegisterReal("compose", NewComposeMarshaler)
	RegisterImage("compose-image", NewComposeImageMarshaler)
	//	RegisterImage("of-gray", NewOfGrayMarshaler)
	//	RegisterImage("of-rgb", NewOfRGBMarshaler)
}

// NewComposeMarshaler returns a Compose transform which can be decoded into.
func NewComposeMarshaler() Real {
	return &Compose{new(RealMarshaler), new(RealMarshaler)}
}

// NewComposeImageMarshaler returns a ComposeImage transform which can be decoded into.
func NewComposeImageMarshaler() Image {
	return &ComposeImage{new(RealMarshaler), new(ImageMarshaler)}
}

//	// NewOfGrayMarshaler returns an OfGray transform which can be decoded into.
//	func NewOfGrayMarshaler() Image { return &OfGray{NewRealMarshaler()} }
//
//	// NewOfRGBMarshaler returns an OfRGB transform which can be decoded into.
//	func NewOfRGBMarshaler() Image { return &OfRGB{NewRealMarshaler()} }

// Compose computes Outer(Inner(x)).
// Compose is itself a Real transform, enabling chains of functions.
type Compose struct {
	Outer, Inner Real
}

func (phi *Compose) Rate() int {
	return phi.Outer.Rate() * phi.Inner.Rate()
}

func (phi *Compose) Apply(x *rimg64.Multi) (*rimg64.Multi, error) {
	z, err := phi.Inner.Apply(x)
	if err != nil {
		return nil, err
	}
	return phi.Outer.Apply(z)
}

// ComposeImage computes Outer(Inner(x)).
// Unlike a Compose transform, the Inner function is computed
// directly on the integer-valued image.
type ComposeImage struct {
	Outer Real
	Inner Image
}

func (phi *ComposeImage) Rate() int {
	return phi.Outer.Rate() * phi.Inner.Rate()
}

func (phi *ComposeImage) Apply(im image.Image) (*rimg64.Multi, error) {
	z, err := phi.Inner.Apply(im)
	if err != nil {
		return nil, err
	}
	return phi.Outer.Apply(z)
}

// Gray describes a real transform applied to the RGB channels of an image.
// Gray{phi} behaves like ComposeImage{phi, NewGray()}.
type OfGray struct{ Real }

func (phi *OfGray) Rate() int { return phi.Real.Rate() }

func (phi *OfGray) Apply(im image.Image) (*rimg64.Multi, error) {
	// toGray never returns an error.
	x, _ := toGray(im)
	return phi.Real.Apply(x)
}

// OfRGB describes a real transform applied to the RGB channels of an image.
// OfRGB{phi} behaves like ComposeImage{phi, NewRGB()}.
type OfRGB struct{ Real }

func (phi *OfRGB) Rate() int { return phi.Real.Rate() }

func (phi *OfRGB) Apply(im image.Image) (*rimg64.Multi, error) {
	return phi.Real.Apply(rimg64.FromColor(im))
}

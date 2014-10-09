package feat

import (
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
)

func init() {
	RegisterReal("compose", NewComposeSpec)
	RegisterImage("compose-image", NewComposeImageSpec)
}

// NewComposeSpec returns a Compose transform which can be decoded into.
func NewComposeSpec() RealSpec {
	return new(composeSpec)
}

// NewComposeImageSpec returns a ComposeImage transform which can be decoded into.
func NewComposeImageSpec() ImageSpec {
	return new(composeImageSpec)
}

// Compose computes Outer(Inner(x)).
// Compose is itself a Real transform, enabling chains of functions.
type Compose struct {
	Outer, Inner RealMarshalable
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

func (phi *Compose) Marshaler() *RealMarshaler {
	// Obtain marshaler for each member.
	spec := &composeSpec{
		Outer: phi.Outer.Marshaler(),
		Inner: phi.Inner.Marshaler(),
	}
	return &RealMarshaler{"compose", spec}
}

// composeSpec contains the contents of Compose.
type composeSpec struct {
	Outer, Inner *RealMarshaler
}

func (m *composeSpec) Transform() RealMarshalable {
	// Obtain transform from each marshaler.
	return &Compose{
		Outer: m.Outer.Spec.Transform(),
		Inner: m.Inner.Spec.Transform(),
	}
}

// ComposeImage computes Outer(Inner(x)).
// Unlike a Compose transform, the Inner function is computed
// directly on the integer-valued image.
type ComposeImage struct {
	Outer RealMarshalable
	Inner ImageMarshalable
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

func (phi *ComposeImage) Marshaler() *ImageMarshaler {
	// Obtain marshaler for each member.
	spec := &composeImageSpec{
		Outer: phi.Outer.Marshaler(),
		Inner: phi.Inner.Marshaler(),
	}
	return &ImageMarshaler{"compose-image", spec}
}

// composeImageSpec contains the contents of Compose.
type composeImageSpec struct {
	Outer *RealMarshaler
	Inner *ImageMarshaler
}

func (m *composeImageSpec) Transform() ImageMarshalable {
	// Obtain transform from each marshaler.
	return &ComposeImage{
		Outer: m.Outer.Spec.Transform(),
		Inner: m.Inner.Spec.Transform(),
	}
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

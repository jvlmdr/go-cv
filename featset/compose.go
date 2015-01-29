package featset

import (
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
)

func init() {
	RegisterReal("compose", newComposeMarshaler)
	RegisterImage("compose-image", newComposeImageMarshaler)
}

func newComposeMarshaler() Real {
	return &Compose{
		Outer: new(RealMarshaler),
		Inner: new(RealMarshaler),
	}
}

func newComposeImageMarshaler() Image {
	return &ComposeImage{
		Outer: new(RealMarshaler),
		Inner: new(ImageMarshaler),
	}
}

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

func (phi *Compose) Size(x image.Point) image.Point {
	return phi.Outer.Size(phi.Inner.Size(x))
}

func (phi *Compose) Channels() int {
	return phi.Outer.Channels()
}

func (phi *Compose) Marshaler() *RealMarshaler {
	// Obtain marshaler for each member.
	return &RealMarshaler{"compose", &Compose{
		Outer: phi.Outer.Marshaler(),
		Inner: phi.Inner.Marshaler(),
	}}
}

func (phi *Compose) Transform() Real {
	return &Compose{
		Outer: phi.Outer.Transform(),
		Inner: phi.Inner.Transform(),
	}
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

func (phi *ComposeImage) Size(x image.Point) image.Point {
	return phi.Outer.Size(phi.Inner.Size(x))
}

func (phi *ComposeImage) Channels() int {
	return phi.Outer.Channels()
}

func (phi *ComposeImage) Marshaler() *ImageMarshaler {
	// Obtain marshaler for each member.
	spec := &ComposeImage{
		Outer: phi.Outer.Marshaler(),
		Inner: phi.Inner.Marshaler(),
	}
	return &ImageMarshaler{"compose-image", spec}
}

func (phi *ComposeImage) Transform() Image {
	return &ComposeImage{
		Outer: phi.Outer.Transform(),
		Inner: phi.Inner.Transform(),
	}
}

// Gray describes a real transform applied to the RGB channels of an image.
// Gray{phi} behaves like ComposeImage{phi, NewGray()}.
type OfGray struct{ Real }

func (phi *OfGray) Rate() int { return phi.Real.Rate() }

func (phi *OfGray) Apply(im image.Image) (*rimg64.Multi, error) {
	return phi.Real.Apply(toGray(im))
}

// OfRGB describes a real transform applied to the RGB channels of an image.
// OfRGB{phi} behaves like ComposeImage{phi, NewRGB()}.
type OfRGB struct{ Real }

func (phi *OfRGB) Rate() int { return phi.Real.Rate() }

func (phi *OfRGB) Apply(im image.Image) (*rimg64.Multi, error) {
	return phi.Real.Apply(rimg64.FromColor(im))
}

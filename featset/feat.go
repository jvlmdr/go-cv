package featset

import "github.com/jvlmdr/go-cv/feat"

type Image interface {
	feat.Image
	Marshaler() *ImageMarshaler
	Transform() Image
}

type Real interface {
	feat.Real
	Marshaler() *RealMarshaler
	Transform() Real
}

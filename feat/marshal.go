package feat

import (
	"encoding/json"
	"fmt"
)

type ImageMarshalable interface {
	Image
	Marshaler() *ImageMarshaler
}

type RealMarshalable interface {
	Real
	Marshaler() *RealMarshaler
}

type ImageSpec interface {
	Transform() ImageMarshalable
}

type RealSpec interface {
	Transform() RealMarshalable
}

// ImageMarshaler uses the default JSON marshaler and defines its own unmarshaler.
type ImageMarshaler struct {
	Name string
	Spec ImageSpec
}

func (m *ImageMarshaler) UnmarshalJSON(data []byte) error {
	var x struct {
		Name string
		Spec json.RawMessage
	}
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if len(x.Name) == 0 {
		return fmt.Errorf("no feature name specified")
	}
	create, ok := ImageTransforms[x.Name]
	if !ok {
		return fmt.Errorf(`unknown feature: "%s"`, x.Name)
	}
	t := create()
	if err := json.Unmarshal(x.Spec, t); err != nil {
		return err
	}
	m.Name = x.Name
	m.Spec = t
	return nil
}

type RealMarshaler struct {
	Name string
	Spec RealSpec
}

func (m *RealMarshaler) UnmarshalJSON(data []byte) error {
	var x struct {
		Name string
		Spec json.RawMessage
	}
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if len(x.Name) == 0 {
		return fmt.Errorf("no feature name specified")
	}
	create, ok := RealTransforms[x.Name]
	if !ok {
		return fmt.Errorf(`unknown feature: "%s"`, x.Name)
	}
	t := create()
	if err := json.Unmarshal(x.Spec, t); err != nil {
		return err
	}
	m.Name = x.Name
	m.Spec = t
	return nil
}

// SimpleImageSpec wraps a simple transform.
type SimpleImageSpec struct {
	Phi ImageMarshalable
}

func (m *SimpleImageSpec) Transform() ImageMarshalable {
	return m.Phi
}

func (m *SimpleImageSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Phi)
}

func (m *SimpleImageSpec) UnmarshalJSON(p []byte) error {
	return json.Unmarshal(p, m.Phi)
}

// SimpleRealSpec wraps a simple transform.
type SimpleRealSpec struct {
	Phi RealMarshalable
}

func (m *SimpleRealSpec) Transform() RealMarshalable {
	return m.Phi
}

func (m *SimpleRealSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Phi)
}

func (m *SimpleRealSpec) UnmarshalJSON(p []byte) error {
	return json.Unmarshal(p, m.Phi)
}

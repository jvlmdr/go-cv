package feat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
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

func NewImageSpec(phi ImageMarshalable) ImageSpec {
	return &simpleImageSpec{phi}
}

// simpleImageSpec wraps a simple transform.
type simpleImageSpec struct {
	Phi ImageMarshalable
}

func (m *simpleImageSpec) Transform() ImageMarshalable {
	return m.Phi
}

func (m *simpleImageSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Phi)
}

func (m *simpleImageSpec) UnmarshalJSON(p []byte) error {
	return json.Unmarshal(p, m.Phi)
}

func NewRealSpec(phi RealMarshalable) RealSpec {
	return &simpleRealSpec{phi}
}

// simpleRealSpec wraps a simple transform.
type simpleRealSpec struct {
	Phi RealMarshalable
}

func (m *simpleRealSpec) Transform() RealMarshalable {
	return m.Phi
}

func (m *simpleRealSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Phi)
}

func (m *simpleRealSpec) UnmarshalJSON(p []byte) error {
	return json.Unmarshal(p, m.Phi)
}

func TestRealMarshaler(phi RealMarshalable) error {
	m := phi.Marshaler()
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(m); err != nil {
		return err
	}
	var u RealMarshaler
	if err := json.NewDecoder(&b).Decode(&u); err != nil {
		return err
	}
	got := u.Spec.Transform()
	if !reflect.DeepEqual(phi, got) {
		err := fmt.Errorf("encode and decode: want %#v, got %#v", phi, got)
		return err
	}
	return nil
}

func TestImageMarshaler(phi ImageMarshalable) error {
	m := phi.Marshaler()
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(m); err != nil {
		return err
	}
	var u ImageMarshaler
	if err := json.NewDecoder(&b).Decode(&u); err != nil {
		return err
	}
	got := u.Spec.Transform()
	if !reflect.DeepEqual(phi, got) {
		err := fmt.Errorf("encode and decode: want %#v, got %#v", phi, got)
		return err
	}
	return nil
}

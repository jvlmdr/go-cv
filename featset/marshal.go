package featset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"reflect"

	"github.com/jvlmdr/go-cv/rimg64"
)

// ImageMarshaler uses the default JSON marshaler and defines its own unmarshaler.
type ImageMarshaler struct {
	Name string
	Spec Image `json:",omitempty"`
}

func (m *ImageMarshaler) Rate() int                  { return m.Spec.Rate() }
func (m *ImageMarshaler) Marshaler() *ImageMarshaler { return m }
func (m *ImageMarshaler) Transform() Image           { return m.Spec.Transform() }

func (m *ImageMarshaler) Apply(im image.Image) (*rimg64.Multi, error) { return m.Spec.Apply(im) }
func (m *ImageMarshaler) Size(x image.Point) image.Point              { return m.Spec.Size(x) }
func (m *ImageMarshaler) Channels() int                               { return m.Spec.Channels() }

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
	if len(x.Spec) > 0 {
		if err := json.Unmarshal(x.Spec, t); err != nil {
			return err
		}
	}
	m.Name = x.Name
	m.Spec = t
	return nil
}

type RealMarshaler struct {
	Name string
	Spec Real `json:",omitempty"`
}

func (m *RealMarshaler) Rate() int                 { return m.Spec.Rate() }
func (m *RealMarshaler) Marshaler() *RealMarshaler { return m }
func (m *RealMarshaler) Transform() Real           { return m.Spec.Transform() }

func (m *RealMarshaler) Apply(f *rimg64.Multi) (*rimg64.Multi, error) { return m.Spec.Apply(f) }
func (m *RealMarshaler) Size(x image.Point) image.Point               { return m.Spec.Size(x) }
func (m *RealMarshaler) Channels() int                                { return m.Spec.Channels() }

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
	if len(x.Spec) > 0 {
		if err := json.Unmarshal(x.Spec, t); err != nil {
			return err
		}
	}
	m.Name = x.Name
	m.Spec = t
	return nil
}

func TestRealMarshaler(phi Real) error {
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

func TestImageMarshaler(phi Image) error {
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

package feat

import (
	"encoding/json"
	"fmt"
)

type ImageMarshaler struct {
	Name  string
	Image `json:"Spec"`
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
	m.Image = t
	return nil
}

type RealMarshaler struct {
	Name string
	Real `json:"Spec"`
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
	m.Real = t
	return nil
}

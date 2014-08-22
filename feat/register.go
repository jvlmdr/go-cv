package feat

import (
	"encoding/json"
	"fmt"
)

var Transforms = make(map[string]func() Transform)

func Register(name string, create func() Transform) {
	Transforms[name] = create
}

type Marshaler struct {
	Name      string
	Transform `json:"Spec"`
}

func (m *Marshaler) UnmarshalJSON(data []byte) error {
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
	create, ok := Transforms[x.Name]
	if !ok {
		return fmt.Errorf(`unknown feature: "%s"`, x.Name)
	}
	t := create()
	if err := json.Unmarshal(x.Spec, t); err != nil {
		return err
	}
	m.Name = x.Name
	m.Transform = t
	return nil
}

package feat

import (
	"encoding/json"
	"fmt"
)

var Transforms = make(map[string]func() Transform)

func Register(name string, create func() Transform) {
	Transforms[name] = create
}

type Unmarshaler struct {
	Name string
	Spec Transform
}

func (m *Unmarshaler) UnmarshalJSON(data []byte) error {
	var x struct {
		Name string
		Spec json.RawMessage
	}
	if err := json.Unmarshal(data, &x); err != nil {
		return err
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
	m.Spec = t
	return nil
}

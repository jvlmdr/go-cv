package feat

var Transforms map[string]func() Transform

func init() {
	Transforms = make(map[string]func() Transform)
}

func Register(name string, create func() Transform) {
	Transforms[name] = create
}

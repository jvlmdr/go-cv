package feat

var Transforms = make(map[string]func() Transform)

func Register(name string, create func() Transform) {
	Transforms[name] = create
}

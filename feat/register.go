package feat

var transforms map[string]func() Transform

func init() {
	transforms = make(map[string]func() Transform)
}

func Register(name string, create func() Transform) {
	transforms[name] = create
}

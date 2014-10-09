package feat

var (
	ImageTransforms = make(map[string]func() ImageSpec)
	RealTransforms  = make(map[string]func() RealSpec)
)

func RegisterImage(name string, create func() ImageSpec) {
	ImageTransforms[name] = create
}

func RegisterReal(name string, create func() RealSpec) {
	RealTransforms[name] = create
}

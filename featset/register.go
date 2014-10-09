package featset

var (
	ImageTransforms = make(map[string]func() Image)
	RealTransforms  = make(map[string]func() Real)
)

func RegisterImage(name string, create func() Image) {
	ImageTransforms[name] = create
}

func RegisterReal(name string, create func() Real) {
	RealTransforms[name] = create
}

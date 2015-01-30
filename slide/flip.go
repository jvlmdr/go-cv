package slide

import "github.com/jvlmdr/go-cv/rimg64"

// Flip mirrors an image in x and y.
func Flip(f *rimg64.Image) *rimg64.Image {
	g := rimg64.New(f.Width, f.Height)
	for i := 0; i < f.Width; i++ {
		for j := 0; j < f.Height; j++ {
			g.Set(f.Width-1-i, f.Height-1-j, f.At(i, j))
		}
	}
	return g
}

// FlipMulti mirrors a multi-channel image in x and y.
func FlipMulti(f *rimg64.Multi) *rimg64.Multi {
	g := rimg64.NewMulti(f.Width, f.Height, f.Channels)
	for i := 0; i < f.Width; i++ {
		for j := 0; j < f.Height; j++ {
			for k := 0; k < f.Channels; k++ {
				g.Set(f.Width-1-i, f.Height-1-j, k, f.At(i, j, k))
			}
		}
	}
	return g
}

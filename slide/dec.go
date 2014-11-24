package slide

import "github.com/jvlmdr/go-cv/rimg64"

// Decimate takes every r-th sample starting at (0, 0).
func Decimate(f *rimg64.Image, r int) *rimg64.Image {
	out := ceilDivPt(f.Size(), r)
	g := rimg64.New(out.X, out.Y)
	for i := 0; i < g.Width; i++ {
		for j := 0; j < g.Height; j++ {
			g.Set(i, j, f.At(r*i, r*j))
		}
	}
	return g
}

// DecimateMulti takes every r-th sample starting at (0, 0).
func DecimateMulti(f *rimg64.Multi, r int) *rimg64.Multi {
	out := ceilDivPt(f.Size(), r)
	g := rimg64.NewMulti(out.X, out.Y, f.Channels)
	for i := 0; i < g.Width; i++ {
		for j := 0; j < g.Height; j++ {
			for k := 0; k < g.Channels; k++ {
				g.Set(i, j, k, f.At(r*i, r*j, k))
			}
		}
	}
	return g
}

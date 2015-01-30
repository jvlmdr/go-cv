package slide_test

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-cv/slide"
)

func sqr(x float64) float64 { return x * x }

func randImage(width, height int) *rimg64.Image {
	f := rimg64.New(width, height)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			f.Set(i, j, rand.NormFloat64())
		}
	}
	return f
}

func randMulti(width, height, channels int) *rimg64.Multi {
	f := rimg64.NewMulti(width, height, channels)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			for k := 0; k < channels; k++ {
				f.Set(i, j, k, rand.NormFloat64())
			}
		}
	}
	return f
}

func randBank(m, n, q int) *slide.Bank {
	g := &slide.Bank{
		Width:   m,
		Height:  n,
		Filters: make([]*rimg64.Image, q),
	}
	for i := range g.Filters {
		g.Filters[i] = randImage(m, n)
	}
	return g
}

func randMultiBank(m, n, p, q int) *slide.MultiBank {
	g := &slide.MultiBank{
		Width:    m,
		Height:   n,
		Channels: p,
		Filters:  make([]*rimg64.Multi, q),
	}
	for i := range g.Filters {
		g.Filters[i] = randMulti(m, n, p)
	}
	return g
}

func errIfNotEqImage(f, g *rimg64.Image, eps float64) error {
	if !f.Size().Eq(g.Size()) {
		return fmt.Errorf("different size: %v, %v", f.Size(), g.Size())
	}
	for i := 0; i < f.Width; i++ {
		for j := 0; j < f.Height; j++ {
			a, b := f.At(i, j), g.At(i, j)
			if math.Abs(a-b) > eps*math.Max(math.Abs(a), math.Abs(b)) {
				return fmt.Errorf("different at x %d, y %d: %g, %g", i, j, a, b)
			}
		}
	}
	return nil
}

func errIfNotEqMulti(f, g *rimg64.Multi, eps float64) error {
	if !f.Size().Eq(g.Size()) {
		return fmt.Errorf("different size: %v, %v", f.Size(), g.Size())
	}
	if f.Channels != g.Channels {
		return fmt.Errorf("different channels: %d, %d", f.Channels, g.Channels)
	}
	for i := 0; i < f.Width; i++ {
		for j := 0; j < f.Height; j++ {
			for k := 0; k < f.Channels; k++ {
				a, b := f.At(i, j, k), g.At(i, j, k)
				if math.Abs(a-b) > eps*math.Max(math.Abs(a), math.Abs(b)) {
					return fmt.Errorf("different at x %d, y %d, c %d: %g, %g", i, j, k, a, b)
				}
			}
		}
	}
	return nil
}

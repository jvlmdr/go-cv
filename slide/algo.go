package slide

import (
	"fmt"

	"github.com/jvlmdr/go-cv/rimg64"
)

type Algo int

const (
	Auto Algo = iota
	Naive
	FFT
	BLAS
)

func CorrAlgo(f, g *rimg64.Image, algo Algo) *rimg64.Image {
	switch algo {
	case Auto:
		return convAuto(f, g, true)
	case Naive:
		return CorrNaive(f, g)
	case FFT:
		return CorrFFT(f, g)
	default:
		panic(fmt.Sprintf("unknown algorithm: %v", algo))
	}
}

func CorrMultiAlgo(f, g *rimg64.Multi, algo Algo) *rimg64.Image {
	switch algo {
	case Auto:
		return convMultiAuto(f, g, true)
	case Naive:
		return CorrMultiNaive(f, g)
	case FFT:
		return CorrMultiFFT(f, g)
	default:
		panic(fmt.Sprintf("unknown algorithm: %v", algo))
	}
}

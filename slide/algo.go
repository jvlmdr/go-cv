package slide

import (
	"fmt"

	"github.com/jvlmdr/go-cv/rimg64"
)

type Algo int

const (
	Naive Algo = iota
	FFT
	BLAS
)

func CorrAlgo(f, g *rimg64.Image, algo Algo) *rimg64.Image {
	switch algo {
	case Naive:
		return CorrNaive(f, g)
	case FFT:
		return CorrFFT(f, g)
	default:
		panic(fmt.Sprintf("unknown algorithm: %g", algo))
	}
}

func CorrMultiAlgo(f, g *rimg64.Multi, algo Algo) *rimg64.Image {
	switch algo {
	case Naive:
		return CorrMultiNaive(f, g)
	case FFT:
		return CorrMultiFFT(f, g)
	default:
		panic(fmt.Sprintf("unknown algorithm: %g", algo))
	}
}

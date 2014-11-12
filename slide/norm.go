package slide

import (
	"image"
	"log"
	"math"

	"github.com/jvlmdr/go-cv/rimg64"
)

// CosCorr computes normalized cross-correlation without mean subtraction.
func CosCorr(f, g *rimg64.Image, algo Algo) *rimg64.Image {
	h := CorrAlgo(f, g, algo)
	if h == nil {
		return h
	}
	gInvNorm := invNorm(g)
	fSqr := square(f)
	fSqrSum := rimg64.CumSum(fSqr)
	// Normalize every element in the output.
	for i := 0; i < h.Width; i++ {
		for j := 0; j < h.Height; j++ {
			rect := image.Rect(i, j, i+g.Width, j+g.Height)
			fInvNorm := rectInvNorm(fSqrSum, rect)
			h.Set(i, j, fInvNorm*gInvNorm*h.At(i, j))
		}
	}
	return h
}

func invNorm(f *rimg64.Image) float64 {
	var norm float64
	for i := 0; i < f.Width; i++ {
		for j := 0; j < f.Height; j++ {
			norm += sqr(f.At(i, j))
		}
	}
	norm = math.Sqrt(norm) // This cannot be negative.
	if norm == 0 {
		return 0
	}
	return 1 / norm
}

func square(f *rimg64.Image) *rimg64.Image {
	g := rimg64.New(f.Width, f.Height)
	for i := 0; i < f.Width; i++ {
		for j := 0; j < f.Height; j++ {
			g.Set(i, j, sqr(f.At(i, j)))
		}
	}
	return g
}

// Assumes that the original image was positive.
// Otherwise it is necessary to compute a sum table from the absolute value of the image.
func rectInvNorm(sqrsum *rimg64.Table, rect image.Rectangle) (invnorm float64) {
	// Relative error for floating-point precision.
	const eps = 1e-9
	sumSqr := sqrsum.Rect(rect)
	sumSqrErr := eps * absTableRect(sqrsum, rect)
	if math.Abs(sumSqr) <= sumSqrErr {
		return 0
	}
	// Negative value will result in NaN when taking square root.
	if sumSqr <= 0 {
		log.Printf("sumSqr:  %10.3e +/- %10.3e", sumSqr, sumSqrErr)
		panic("square norm not greater than zero")
	}
	// Can compute inverse norm.
	return 1 / math.Sqrt(sumSqr)
}

// NormCorr computes normalized cross-correlation.
func NormCorrMulti(f, g *rimg64.Image) *rimg64.Image {
	return convAuto(f, g, true)
}

// Used for floating-point precision.
// Rather than A - B - C + D, computes A + B + C + D.
// Only useful when original image was non-negative.
func absTableRect(t *rimg64.Table, r image.Rectangle) float64 {
	s := (*rimg64.Image)(t)
	bnds := image.Rect(0, 0, s.Width, s.Height)
	if !r.In(bnds) {
		panic("out of bounds")
	}
	if r.Dx()*r.Dy() == 0 {
		return 0
	}
	area := s.At(r.Max.X-1, r.Max.Y-1)
	if r.Min.X > 0 {
		// Change from plus to minus.
		area += s.At(r.Min.X-1, r.Max.Y-1)
	}
	if r.Min.Y > 0 {
		// Change from plus to minus.
		area += s.At(r.Max.X-1, r.Min.Y-1)
	}
	if r.Min.X > 0 && r.Min.Y > 0 {
		area += s.At(r.Min.X-1, r.Min.Y-1)
	}
	return area
}

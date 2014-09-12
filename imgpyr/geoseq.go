package imgpyr

import (
	"fmt"
	"math"
)

// Describes a finite geometric sequence.
type GeoSeq struct {
	Start float64
	Step  float64
	Len   int
}

// Generates a sequence from start to the last element up to and including lim.
//
// If step > 1, then lim must be greater than start.
// If step < 1, then lim must be less than start.
func Sequence(start, step, lim float64) GeoSeq {
	if step == 1 {
		panic("step must not be 1")
	}
	if step <= 0 {
		panic("step must be positive")
	}
	// Find n such that lim * step^(n-1) == lim.
	n := math.Log(lim/start)/math.Log(step) + 1
	// Round down.
	m := int(math.Floor(n))
	return GeoSeq{start, step, m}
}

// Generates a sequence from first to last containing n elements.
func LogRange(first, last float64, n int) GeoSeq {
	step := math.Exp(math.Log(last/first) / float64(n-1))
	return GeoSeq{first, step, n}
}

// Returns the i-th value of the progression.
func (seq GeoSeq) At(i int) float64 {
	if i < 0 || i >= seq.Len {
		panic(fmt.Sprintf("out of range: %d", i))
	}
	return seq.Start * math.Pow(seq.Step, float64(i))
}

// Returns the (floating point) index of the x in the progression.
func (seq GeoSeq) Inv(x float64) float64 {
	return math.Log(x/seq.Start) / math.Log(seq.Step)
}

// Returns a reversed sequence.
func (seq GeoSeq) Reverse() GeoSeq {
	return GeoSeq{seq.At(seq.Len - 1), 1 / seq.Step, seq.Len}
}

func (seq GeoSeq) Elems() []float64 {
	x := make([]float64, seq.Len)
	for i := 0; i < seq.Len; i++ {
		x[i] = seq.At(i)
	}
	return x
}

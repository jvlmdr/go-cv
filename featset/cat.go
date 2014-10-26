package featset

import (
	"errors"
	"fmt"
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
)

func init() {
	RegisterReal("concat", newConcatMarshaler)
}

type RealList interface {
	Len() int
	At(int) Real
}

type RealSlice []Real

func (fs RealSlice) Len() int      { return len(fs) }
func (fs RealSlice) At(i int) Real { return fs[i] }

// Need a concrete slice to unmarshal unknown number.
// Could use a linked list instead but it might be messier.
type realMarshalerSlice []*RealMarshaler

func (fs realMarshalerSlice) Len() int      { return len(fs) }
func (fs realMarshalerSlice) At(i int) Real { return fs[i] }

func newConcatMarshaler() Real {
	return &Concat{realMarshalerSlice{}}
}

// Concat concatenates the channels of a collection of transforms.
// Each transform must output the same size and have the same rate.
type Concat struct {
	Elems RealList
}

func (phi *Concat) Rate() int {
	var rate int
	for i := 0; i < phi.Elems.Len(); i++ {
		elem := phi.Elems.At(i)
		if i == 0 {
			rate = elem.Rate()
			continue
		}
		if elem.Rate() != rate {
			panic("transforms have different rates")
		}
	}
	return rate
}

func (phi *Concat) Apply(x *rimg64.Multi) (*rimg64.Multi, error) {
	ys := make([]*rimg64.Multi, phi.Elems.Len())
	for i := 0; i < phi.Elems.Len(); i++ {
		elem := phi.Elems.At(i)
		// Execute each transform.
		y, err := elem.Apply(x)
		if err != nil {
			return nil, err
		}
		ys[i] = y
	}
	// If all nil, then return nil.
	allNil := true
	for _, y := range ys {
		if y != nil {
			allNil = false
			break
		}
	}
	if allNil {
		return nil, nil
	}
	// Check that sizes are all the same.
	var size image.Point
	for i, y := range ys {
		if y == nil {
			return nil, errors.New("some transforms are nil")
		}
		if i == 0 {
			size = y.Size()
			continue
		}
		if !y.Size().Eq(size) {
			return nil, fmt.Errorf("different image sizes: %v, %v", size, y.Size())
		}
	}
	var channels int
	for _, y := range ys {
		channels += y.Channels
	}
	// Copy into one image.
	z := rimg64.NewMulti(size.X, size.Y, channels)
	q := 0
	for _, y := range ys {
		for u := 0; u < size.X; u++ {
			for v := 0; v < size.Y; v++ {
				for p := 0; p < y.Channels; p++ {
					z.Set(u, v, q+p, y.At(u, v, p))
				}
			}
		}
		q += y.Channels
	}
	return z, nil
}

func (phi *Concat) Marshaler() *RealMarshaler {
	// Obtain marshaler for each member.
	ms := make([]*RealMarshaler, phi.Elems.Len())
	for i := 0; i < phi.Elems.Len(); i++ {
		ms[i] = phi.Elems.At(i).Marshaler()
	}
	return &RealMarshaler{"concat", &Concat{realMarshalerSlice(ms)}}
}

func (phi *Concat) Transform() Real {
	// Obtain marshaler for each member.
	fs := make([]Real, phi.Elems.Len())
	for i := 0; i < phi.Elems.Len(); i++ {
		fs[i] = phi.Elems.At(i).Transform()
	}
	return &Concat{RealSlice(fs)}
}

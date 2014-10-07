package convfeat

import (
	"image"

	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-cv/slide"
)

// FilterBank describes a collection of single-channel filters.
// All filters must have the same dimension.
type FilterBank struct {
	Field image.Point
	List  []*rimg64.Image
}

func (bank *FilterBank) Corr(x *rimg64.Image) *rimg64.Multi {
	size := slide.ValidSize(x.Size(), bank.Field)
	y := rimg64.NewMulti(size.X, size.Y, len(bank.List))
	// Convolve y with each filter in the list.
	for i, a := range bank.List {
		y.SetChannel(i, slide.Corr(x, a))
	}
	return y
}

// FilterBankMulti describes a collection of multi-channel filters.
// All filters must have the same dimension.
type FilterBankMulti struct {
	Field image.Point
	NumIn int
	List  []*rimg64.Multi
}

func (bank *FilterBankMulti) Corr(x *rimg64.Multi) *rimg64.Multi {
	size := slide.ValidSize(x.Size(), bank.Field)
	y := rimg64.NewMulti(size.X, size.Y, len(bank.List))
	// Convolve y with each filter in the list.
	for i, a := range bank.List {
		y.SetChannel(i, slide.CorrMulti(x, a))
	}
	return y
}

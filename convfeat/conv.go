package convfeat

import (
	"fmt"

	"github.com/jvlmdr/go-cv/feat"
	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-cv/slide"
)

func init() {
	feat.RegisterReal("conv", func() feat.Real { return new(ConvMulti) })
	feat.RegisterReal("conv-each", func() feat.Real { return new(ConvEach) })
}

// ConvMulti represents multi-channel convolution.
type ConvMulti struct {
	Filters *FilterBankMulti
}

func (phi *ConvMulti) Rate() int { return 1 }

func (phi *ConvMulti) Apply(x *rimg64.Multi) (*rimg64.Multi, error) {
	if x.Channels != phi.Filters.NumIn {
		err := fmt.Errorf(
			"channels: image has %d, filter bank has %d",
			x.Channels, phi.Filters.NumIn,
		)
		return nil, err
	}
	return phi.Filters.Corr(x), nil
}

// ConvEach applies the same single-channel filters to every channel.
type ConvEach struct {
	Filters *FilterBank
}

func (phi *ConvEach) Rate() int { return 1 }

func (phi *ConvEach) Apply(x *rimg64.Multi) (*rimg64.Multi, error) {
	channels := x.Channels * len(phi.Filters.List)
	size := slide.ValidSize(x.Size(), phi.Filters.Field)
	y := rimg64.NewMulti(size.X, size.Y, channels)
	var n int
	for i := 0; i < x.Channels; i++ {
		// Convolve each channel of the input with the bank.
		yi := phi.Filters.Corr(x.Channel(i))
		for j := 0; j < yi.Channels; j++ {
			// Copy the channels into the output.
			y.SetChannel(n, yi.Channel(j))
			n++
		}
	}
	return y, nil
}

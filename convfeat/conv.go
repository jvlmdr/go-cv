package convfeat

import (
	"fmt"

	"github.com/jvlmdr/go-cv/featset"
	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-cv/slide"
)

func init() {
	featset.RegisterReal("conv", func() featset.Real { return new(ConvMulti) })
	featset.RegisterReal("conv-each", func() featset.Real { return new(ConvEach) })
	featset.RegisterReal("add-const", func() featset.Real { return new(AddConst) })
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

func (phi *ConvMulti) Marshaler() *featset.RealMarshaler {
	return &featset.RealMarshaler{"conv", phi}
}

func (phi *ConvMulti) Transform() featset.Real { return phi }

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

func (phi *ConvEach) Marshaler() *featset.RealMarshaler {
	return &featset.RealMarshaler{"conv-each", phi}
}

func (phi *ConvEach) Transform() featset.Real { return phi }

// AddConst adds a constant to every pixel.
type AddConst []float64

func (phi *AddConst) Rate() int { return 1 }

func (phi *AddConst) Apply(x *rimg64.Multi) (*rimg64.Multi, error) {
	if x.Channels != len(*phi) {
		err := fmt.Errorf("channels: image has %d, filter bank has %d", x.Channels, len(*phi))
		return nil, err
	}
	y := rimg64.NewMulti(x.Width, x.Height, x.Channels)
	for u := 0; u < x.Width; u++ {
		for v := 0; v < x.Height; v++ {
			for p := 0; p < x.Channels; p++ {
				y.Set(u, v, p, x.At(u, v, p)+(*phi)[p])
			}
		}
	}
	return y, nil
}

func (phi *AddConst) Marshaler() *featset.RealMarshaler {
	return &featset.RealMarshaler{"add-const", phi}
}

func (phi *AddConst) Transform() featset.Real { return phi }

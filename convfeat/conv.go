package convfeat

import (
	"fmt"
	"image"

	"github.com/jvlmdr/go-cv/featset"
	"github.com/jvlmdr/go-cv/rimg64"
	"github.com/jvlmdr/go-cv/slide"
)

func init() {
	featset.RegisterReal("conv", func() featset.Real { return new(ConvMulti) })
	featset.RegisterReal("conv-each", func() featset.Real { return new(ConvEach) })
	featset.RegisterReal("add-const", func() featset.Real { return new(AddConst) })
	featset.RegisterReal("scale", func() featset.Real { return new(Scale) })
}

// ConvMulti represents multi-channel convolution.
type ConvMulti struct {
	Stride  int
	Filters *slide.MultiBank
}

func (phi *ConvMulti) Rate() int { return phi.Stride }

func (phi *ConvMulti) Apply(x *rimg64.Multi) (*rimg64.Multi, error) {
	if x.Channels != phi.Filters.Channels {
		err := fmt.Errorf(
			"channels: image has %d, filter bank has %d",
			x.Channels, phi.Filters.Channels,
		)
		return nil, err
	}
	if phi.Stride <= 1 {
		return slide.CorrMultiBankBLAS(x, phi.Filters)
	}
	return slide.CorrMultiBankStrideBLAS(x, phi.Filters, phi.Stride)
}

func (phi *ConvMulti) Marshaler() *featset.RealMarshaler {
	return &featset.RealMarshaler{"conv", phi}
}

func (phi *ConvMulti) Transform() featset.Real { return phi }

// ConvEach applies the same single-channel filters to every channel.
type ConvEach struct {
	Filters *slide.Bank
}

func (phi *ConvEach) Rate() int { return 1 }

func (phi *ConvEach) Apply(x *rimg64.Multi) (*rimg64.Multi, error) {
	channels := x.Channels * len(phi.Filters.Filters)
	field := image.Pt(phi.Filters.Width, phi.Filters.Height)
	size := slide.ValidSize(x.Size(), field)
	y := rimg64.NewMulti(size.X, size.Y, channels)
	var n int
	for i := 0; i < x.Channels; i++ {
		// Convolve each channel of the input with the bank.
		yi, err := slide.CorrBankBLAS(x.Channel(i), phi.Filters)
		if err != nil {
			return nil, err
		}
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

// Scale multiplies every pixel by a constant.
type Scale float64

func (phi *Scale) Rate() int { return 1 }

func (phi *Scale) Apply(x *rimg64.Multi) (*rimg64.Multi, error) {
	y := rimg64.NewMulti(x.Width, x.Height, x.Channels)
	for u := 0; u < x.Width; u++ {
		for v := 0; v < x.Height; v++ {
			for p := 0; p < x.Channels; p++ {
				y.Set(u, v, p, float64(*phi)*x.At(u, v, p))
			}
		}
	}
	return y, nil
}

func (phi *Scale) Marshaler() *featset.RealMarshaler {
	return &featset.RealMarshaler{"scale", phi}
}

func (phi *Scale) Transform() featset.Real { return phi }

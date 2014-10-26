package featset

import "github.com/jvlmdr/go-cv/rimg64"

func init() {
	RegisterReal("channel-interval", func() Real { return new(ChannelInterval) })
	RegisterReal("select-channels", func() Real { return new(SelectChannels) })
}

// ChannelInterval selects channels in [a, b).
type ChannelInterval struct{ A, B int }

func (phi *ChannelInterval) Rate() int { return 1 }

func (phi *ChannelInterval) Apply(f *rimg64.Multi) (*rimg64.Multi, error) {
	g := rimg64.NewMulti(f.Width, f.Height, phi.B-phi.A)
	for u := 0; u < f.Width; u++ {
		for v := 0; v < f.Height; v++ {
			for p := phi.A; p < phi.B; p++ {
				g.Set(u, v, p-phi.A, f.At(u, v, p))
			}
		}
	}
	return g, nil
}

func (phi *ChannelInterval) Marshaler() *RealMarshaler {
	return &RealMarshaler{"channel-interval", phi}
}

func (phi *ChannelInterval) Transform() Real { return phi }

// SelectChannels takes a subset of channels in the given order.
type SelectChannels struct {
	Channels []int
}

func (phi *SelectChannels) Rate() int { return 1 }

func (phi *SelectChannels) Apply(f *rimg64.Multi) (*rimg64.Multi, error) {
	g := rimg64.NewMulti(f.Width, f.Height, len(phi.Channels))
	for u := 0; u < f.Width; u++ {
		for v := 0; v < f.Height; v++ {
			for i, p := range phi.Channels {
				g.Set(u, v, i, f.At(u, v, p))
			}
		}
	}
	return g, nil
}

func (phi *SelectChannels) Marshaler() *RealMarshaler {
	return &RealMarshaler{"select-channels", phi}
}

func (phi *SelectChannels) Transform() Real { return phi }

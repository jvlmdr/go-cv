package hog

import (
	"image"
	"math"

	"github.com/jvlmdr/go-cv/feat"
	"github.com/jvlmdr/go-cv/rimg64"
)

func init() {
	feat.RegisterImage("hog", func() feat.ImageSpec {
		return feat.NewImageSpec(new(Transform))
	})
}

type Transform struct {
	Conf Config
}

func (t Transform) Rate() int {
	// For now, down-sample rate equals size of cell.
	return t.Conf.CellSize
}

func (t Transform) Apply(im image.Image) (*rimg64.Multi, error) {
	return HOG(rimg64.FromColor(im), t.Conf), nil
}

func (t Transform) Marshaler() *feat.ImageMarshaler {
	return &feat.ImageMarshaler{"hog", feat.NewImageSpec(new(Transform))}
}

func FGMRConfig(sbin int) Config {
	return Config{
		Angles:   9,
		CellSize: sbin,
	}
}

type Config struct {
	// Number of discrete orientations.
	Angles int
	// Number of pixels to one side of a square cell.
	CellSize int

	// Do not include contrast-sensitive features.
	NoContrastVar bool
	// Do not include contrast-insensitive features.
	NoContrastInvar bool
	// Do not include texture features.
	NoTexture bool
	// Do not let gradient intensities be less than some value.
	NoClip bool
}

type point struct {
	X, Y float64
}

// Returns gradient with greatest magnitude across all channels.
// 1 <= x <= width-2, 1 <= y <= height-2
func maxGrad(f *rimg64.Multi, x, y int) (point, float64) {
	var (
		grad point
		max  float64
	)
	for d := 0; d < f.Channels; d++ {
		p := point{
			f.At(x+1, y, d) - f.At(x-1, y, d),
			f.At(x, y+1, d) - f.At(x, y-1, d),
		}
		v := p.X*p.X + p.Y*p.Y
		if v > max {
			grad, max = p, v
		}
	}
	return grad, max
}

// Returns an index in {0, ..., 2*n-1}.
func quantAngle(grad point, n int) int {
	var (
		q   int     = 0
		max float64 = 0
	)
	for i := 0; i < n; i++ {
		theta := float64(i) / float64(n) * math.Pi
		dot := grad.X*math.Cos(theta) + grad.Y*math.Sin(theta)
		if dot > max {
			q, max = i, dot
		} else if -dot > max {
			q, max = i+n, -dot
		}
	}
	return q
}

func adjSum(f *rimg64.Image, x1, y1, x2, y2 int) float64 {
	return f.At(x1, y1) + f.At(x1, y2) + f.At(x2, y1) + f.At(x2, y2)
}

func HOG(f *rimg64.Multi, conf Config) *rimg64.Multi {
	const eps = 0.0001
	var channels int
	if !conf.NoContrastVar {
		channels += 2 * conf.Angles
	}
	if !conf.NoContrastInvar {
		channels += conf.Angles
	}
	if !conf.NoTexture {
		channels += 4
	}

	// Leave a one-pixel border to compute derivatives.
	inside := image.Rectangle{image.ZP, f.Size()}.Inset(1)
	// Leave a half-cell border.
	half := conf.CellSize / 2
	valid := inside.Inset(half)
	// Number of whole cells inside valid region.
	cells := valid.Size().Div(conf.CellSize)
	if cells.X <= 0 || cells.Y <= 0 {
		return nil
	}
	// Remove one cell on all sides for output.
	out := cells.Sub(image.Pt(2, 2))
	// Region to iterate over.
	size := cells.Mul(conf.CellSize).Add(image.Pt(2*half, 2*half))
	vis := image.Rectangle{inside.Min, inside.Min.Add(size)}

	// Accumulate edges into cell histograms.
	hist := rimg64.NewMulti(cells.X, cells.Y, 2*conf.Angles)
	for a := vis.Min.X; a < vis.Max.X; a++ {
		for b := vis.Min.Y; b < vis.Max.Y; b++ {
			x, y := a-half-vis.Min.X, b-half-vis.Min.Y
			// Pick channel with strongest gradient.
			grad, v := maxGrad(f, a, b)
			v = math.Sqrt(v)
			// Snap to orientation.
			q := quantAngle(grad, conf.Angles)

			// Add to 4 histograms around pixel using bilinear interpolation.
			xp := (float64(x)+0.5)/float64(conf.CellSize) - 0.5
			yp := (float64(y)+0.5)/float64(conf.CellSize) - 0.5
			// Extract integer and fractional part.
			ixp, vx0 := modf(xp)
			iyp, vy0 := modf(yp)
			// Complement of fraction part.
			vx1 := 1 - vx0
			vy1 := 1 - vy0

			if ixp >= 0 && iyp >= 0 {
				addToMulti(hist, ixp, iyp, q, vx1*vy1*v)
			}
			if ixp+1 < cells.X && iyp >= 0 {
				addToMulti(hist, ixp+1, iyp, q, vx0*vy1*v)
			}
			if ixp >= 0 && iyp+1 < cells.Y {
				addToMulti(hist, ixp, iyp+1, q, vx1*vy0*v)
			}
			if ixp+1 < cells.X && iyp+1 < cells.Y {
				addToMulti(hist, ixp+1, iyp+1, q, vx0*vy0*v)
			}
		}
	}

	// compute energy in each block by summing over orientations
	norm := rimg64.New(cells.X, cells.Y)
	for x := 0; x < cells.X; x++ {
		for y := 0; y < cells.Y; y++ {
			for d := 0; d < conf.Angles; d++ {
				s := hist.At(x, y, d) + hist.At(x, y, d+conf.Angles)
				addTo(norm, x, y, s*s)
			}
		}
	}

	feat := rimg64.NewMulti(out.X, out.Y, channels)
	for x := 0; x < out.X; x++ {
		for y := 0; y < out.Y; y++ {
			a, b := x+1, y+1
			// Normalization factors.
			var n [4]float64
			n[0] = 1 / math.Sqrt(adjSum(norm, a, b, a+1, b+1)+eps)
			n[1] = 1 / math.Sqrt(adjSum(norm, a, b, a+1, b-1)+eps)
			n[2] = 1 / math.Sqrt(adjSum(norm, a, b, a-1, b+1)+eps)
			n[3] = 1 / math.Sqrt(adjSum(norm, a, b, a-1, b-1)+eps)
			var off int

			// Contrast-sensitive features.
			if !conf.NoContrastVar {
				for d := 0; d < 2*conf.Angles; d++ {
					h := hist.At(a, b, d)
					var sum float64
					for _, ni := range n {
						val := h * ni
						if !conf.NoClip {
							val = math.Min(val, 0.2)
						}
						sum += val
					}
					feat.Set(x, y, off+d, sum/2)
				}
				off += 2 * conf.Angles
			}

			// Contrast-insensitive features.
			if !conf.NoContrastInvar {
				for d := 0; d < conf.Angles; d++ {
					h := hist.At(a, b, d) + hist.At(a, b, conf.Angles+d)
					var sum float64
					for _, ni := range n {
						val := h * ni
						if !conf.NoClip {
							val = math.Min(val, 0.2)
						}
						sum += val
					}
					feat.Set(x, y, off+d, sum/2)
				}
				off += conf.Angles
			}

			// Texture features.
			if !conf.NoTexture {
				for i, ni := range n {
					var sum float64
					for d := 0; d < 2*conf.Angles; d++ {
						h := hist.At(a, b, d)
						val := h * ni
						if !conf.NoClip {
							val = math.Min(val, 0.2)
						}
						sum += val
					}
					feat.Set(x, y, off+i, sum/math.Sqrt(float64(2*conf.Angles)))
				}
				off += 4
			}
		}
	}
	return feat
}

func FeatSize(pix image.Point, conf Config) image.Point {
	// Leave a one-pixel border to compute derivatives.
	inside := image.Rectangle{image.ZP, pix}.Inset(1)
	// Leave a half-cell border.
	half := conf.CellSize / 2
	valid := inside.Inset(half)
	// Number of whole cells inside valid region.
	cells := valid.Size().Div(conf.CellSize)
	if cells.X <= 0 || cells.Y <= 0 {
		return image.ZP
	}
	// Remove one cell on all sides for output.
	out := cells.Sub(image.Pt(2, 2))
	return out
}

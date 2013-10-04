package hog

import (
	"github.com/jackvalmadre/go-cv/rimg64"

	"image"
	"math"
)

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
}

type point struct {
	X, Y float64
}

// Returns gradient with greatest magnitude across all channels.
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

func adjSum(f *rimg64.Image, x, y int) float64 {
	return f.At(x, y) + f.At(x, y+1) + f.At(x+1, y) + f.At(x+1, y+1)
}

func HOG(f *rimg64.Multi, conf Config) *rimg64.Multi {
	const eps = 0.0001

	// Number of cells.
	cells := image.Pt(
		int(math.Floor(float64(f.Width)/float64(conf.CellSize))),
		int(math.Floor(float64(f.Height)/float64(conf.CellSize))),
	)
	// Pixels which are covered by cells.
	visible := cells.Mul(conf.CellSize)
	// Size of output image.
	// Exclude edge cells.
	out := image.Pt(max(cells.X-2, 0), max(cells.Y-2, 0))
	channels := 3*conf.Angles + 4

	// Accumulate edges into cell histograms.
	hist := rimg64.NewMulti(cells.X, cells.Y, 2*conf.Angles)
	for x := 1; x < visible.X-1; x++ {
		for y := 1; y < visible.Y-1; y++ {
			// Pick channel with strongest gradient.
			grad, v := maxGrad(f, x, y)
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
			n[0] = 1 / math.Sqrt(adjSum(norm, a, b)+eps)
			n[1] = 1 / math.Sqrt(adjSum(norm, a, b-1)+eps)
			n[2] = 1 / math.Sqrt(adjSum(norm, a-1, b)+eps)
			n[3] = 1 / math.Sqrt(adjSum(norm, a-1, b-1)+eps)

			// Directed edges.
			for d := 0; d < 2*conf.Angles; d++ {
				h := hist.At(a, b, d)
				var sum float64
				for _, ni := range n {
					sum += math.Min(h*ni, 0.2)
				}
				feat.Set(x, y, d, sum/2)
			}

			// Un-directed edges.
			off := 2 * conf.Angles
			for d := 0; d < conf.Angles; d++ {
				h := hist.At(a, b, d) + hist.At(a, b, conf.Angles+d)
				var sum float64
				for _, ni := range n {
					sum += math.Min(h*ni, 0.2)
				}
				feat.Set(x, y, off+d, sum/2)
			}

			// Texture features.
			off = 3 * conf.Angles
			for i, ni := range n {
				var sum float64
				for d := 0; d < 2*conf.Angles; d++ {
					h := hist.At(a, b, d)
					sum += math.Min(h*ni, 0.2)
				}
				feat.Set(x, y, off+i, 0.2357*sum)
			}
		}
	}

	return feat
}

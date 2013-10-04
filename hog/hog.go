package hog

import (
	"github.com/jackvalmadre/go-cv/rimg64"

	"math"
)

const NumAngles = 9

func HOG(f *rimg64.Multi, sbin int) *rimg64.Multi {
	const eps = 0.0001

	//	var uu [NumAngles]float64
	//	var vv [NumAngles]float64
	//	for i := range uu {
	//		uu[i] = math.Cos(float64(i)/NumAngles*math.Pi)
	//		vv[i] = math.Sin(float64(i)/NumAngles*math.Pi)
	//	}

	uu := [NumAngles]float64{1.0000,
		0.9397,
		0.7660,
		0.500,
		0.1736,
		-0.1736,
		-0.5000,
		-0.7660,
		-0.9397,
	}
	vv := [NumAngles]float64{0.0000,
		0.3420,
		0.6428,
		0.8660,
		0.9848,
		0.9848,
		0.8660,
		0.6428,
		0.3420,
	}

	dims := [2]int{f.Height, f.Width}
	cells := [2]int{
		round(float64(dims[0]) / float64(sbin)),
		round(float64(dims[1]) / float64(sbin)),
	}
	out := [3]int{
		max(cells[0]-2, 0),
		max(cells[1]-2, 0),
		27 + 4,
	}

	hist := rimg64.NewMulti(cells[1], cells[0], 2*NumAngles)
	norm := rimg64.New(cells[1], cells[0])
	feat := rimg64.NewMulti(out[1], out[0], 3*NumAngles+4)

	visible := [2]int{cells[0] * sbin, cells[1] * sbin}

	for x := 1; x < visible[1]-1; x++ {
		for y := 1; y < visible[0]-1; y++ {
			a := min(x, dims[1]-2)
			b := min(y, dims[0]-2)

			// pick channel with strongest gradient
			var (
				dx, dy, v float64
			)
			for d := 0; d < f.Channels; d++ {
				p := f.At(a+1, b, d) - f.At(a-1, b, d)
				q := f.At(a, b+1, d) - f.At(a, b-1, d)
				r := p*p + q*q
				if r > v {
					dx, dy, v = p, q, r
				}
			}

			// snap to one of 18 orientations
			var (
				best_dot float64 = 0
				best_o   int     = 0
			)
			for o := 0; o < NumAngles; o++ {
				dot := uu[o]*dx + vv[o]*dy
				if dot > best_dot {
					best_dot, best_o = dot, o
				} else if -dot > best_dot {
					best_dot, best_o = -dot, o+NumAngles
				}
			}

			// add to 4 histograms around pixel using bilinear interpolation
			xp := (float64(x)+0.5)/float64(sbin) - 0.5
			yp := (float64(y)+0.5)/float64(sbin) - 0.5
			ixp, vx0 := modf(xp)
			iyp, vy0 := modf(yp)
			vx1 := 1 - vx0
			vy1 := 1 - vy0
			v = math.Sqrt(v)

			if ixp >= 0 && iyp >= 0 {
				addToMulti(hist, ixp, iyp, best_o, vx1*vy1*v)
			}
			if ixp+1 < cells[1] && iyp >= 0 {
				addToMulti(hist, ixp+1, iyp, best_o, vx0*vy1*v)
			}
			if ixp >= 0 && iyp+1 < cells[0] {
				addToMulti(hist, ixp, iyp+1, best_o, vx1*vy0*v)
			}
			if ixp+1 < cells[1] && iyp+1 < cells[0] {
				addToMulti(hist, ixp+1, iyp+1, best_o, vx0*vy0*v)
			}
		}
	}

	// compute energy in each block by summing over orientations
	for x := 0; x < cells[1]; x++ {
		for y := 0; y < cells[0]; y++ {
			for o := 0; o < 9; o++ {
				s := hist.At(x, y, o) + hist.At(x, y, o+NumAngles)
				addTo(norm, x, y, s*s)
			}
		}
	}

	for x := 0; x < out[1]; x++ {
		for y := 0; y < out[0]; y++ {
			var a, b int
			a, b = x+1, y+1
			n1 := 1 / math.Sqrt(norm.At(a, b)+norm.At(a, b+1)+norm.At(a+1, b)+norm.At(a+1, b+1)+eps)
			a, b = x+1, y
			n2 := 1 / math.Sqrt(norm.At(a, b)+norm.At(a, b+1)+norm.At(a+1, b)+norm.At(a+1, b+1)+eps)
			a, b = x, y+1
			n3 := 1 / math.Sqrt(norm.At(a, b)+norm.At(a, b+1)+norm.At(a+1, b)+norm.At(a+1, b+1)+eps)
			a, b = x, y
			n4 := 1 / math.Sqrt(norm.At(a, b)+norm.At(a, b+1)+norm.At(a+1, b)+norm.At(a+1, b+1)+eps)

			// contrast-sensitive features
			var t1, t2, t3, t4 float64
			for o := 0; o < 2*NumAngles; o++ {
				h := hist.At(x+1, y+1, o)
				h1 := math.Min(h*n1, 0.2)
				h2 := math.Min(h*n2, 0.2)
				h3 := math.Min(h*n3, 0.2)
				h4 := math.Min(h*n4, 0.2)
				feat.Set(x, y, o, (h1+h2+h3+h4)/2)
				t1 += h1
				t2 += h2
				t3 += h3
				t4 += h4
			}

			// contrast-insensitive features
			off := 2 * NumAngles
			for o := 0; o < NumAngles; o++ {
				h := hist.At(x+1, y+1, o) + hist.At(x+1, y+1, NumAngles+o)
				h1 := math.Min(h*n1, 0.2)
				h2 := math.Min(h*n2, 0.2)
				h3 := math.Min(h*n3, 0.2)
				h4 := math.Min(h*n4, 0.2)
				feat.Set(x, y, off+o, (h1+h2+h3+h4)/2)
			}

			// texture features
			off = 3 * NumAngles
			feat.Set(x, y, off, 0.2357*t1)
			feat.Set(x, y, off+1, 0.2357*t2)
			feat.Set(x, y, off+2, 0.2357*t3)
			feat.Set(x, y, off+3, 0.2357*t4)
		}
	}

	return feat
}

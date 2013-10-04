package hog

import (
	"code.google.com/p/draw2d/draw2d"
	"github.com/jackvalmadre/go-cv/rimg64"
	"github.com/jackvalmadre/lin-go/vec"

	"image"
	"image/color"
	"image/draw"
	"math"
)

type WeightSet int

const (
	Signed WeightSet = iota
	Pos
	Neg
	Abs
)

func Vis(feat *rimg64.Multi, weights WeightSet, cell int) image.Image {
	if feat.Channels == 31 {
		return Vis(compress(feat, weights), weights, cell)
	}

	width, height := feat.Width*cell, feat.Height*cell
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill background.
	bg := color.Gray{0}
	if weights == Signed {
		bg = color.Gray{128}
	}
	draw.Draw(img, img.Bounds(), image.NewUniform(bg), image.ZP, draw.Src)

	// Rescale intensities to [0, 1].
	var dst vec.Mutable = vec.Slice(feat.Elems)
	var src vec.Const = dst

	if weights == Signed {
		max, _ := vec.Max(vec.Abs(src))
		rescale := func(x float64) float64 {
			return (1 + x/max) / 2
		}
		vec.Copy(dst, vec.Map(src, rescale))
	} else {
		switch weights {
		case Neg:
			src = vec.Scale(-1, src)
		case Abs:
			src = vec.Abs(src)
		default:
		}

		max, _ := vec.Max(src)
		if max <= 0 {
			vec.Copy(dst, vec.Zeros(src.Len()))
		} else {
			rescale := func(x float64) float64 {
				return math.Max(0, x/max)
			}
			vec.Copy(dst, vec.Map(src, rescale))
		}
	}

	gc := draw2d.NewGraphicContext(img)
	gc.SetLineWidth(1)

	// Draw cells.
	for x := 0; x < feat.Width; x++ {
		for y := 0; y < feat.Height; y++ {
			drawCell(feat, x, y, gc, cell)
		}
	}
	return img
}

// Flattens 31 (or 27) channels down to 9 for visualization.
func compress(src *rimg64.Multi, weights WeightSet) *rimg64.Multi {
	dst := rimg64.NewMulti(src.Width, src.Height, 9)
	for i := 0; i < 27; i++ {
		for x := 0; x < src.Width; x++ {
			for y := 0; y < src.Height; y++ {
				v := src.At(x, y, i)
				switch weights {
				default:
				case Pos:
					v = math.Max(0, v)
				case Neg:
					v = math.Min(0, v)
				case Abs:
					v = math.Abs(v)
				}
				dst.Set(x, y, i%9, dst.At(x, y, i%9)+v)
			}
		}
	}
	return dst
}

func drawCell(feat *rimg64.Multi, i, j int, gc *draw2d.ImageGraphicContext, cell int) {
	u := (float64(i) + 0.5) * float64(cell)
	v := (float64(j) + 0.5) * float64(cell)
	r := float64(cell) / 2

	for k := 0; k < Orientations; k++ {
		x := feat.At(i, j, k)
		x = math.Max(x, 0)
		x = math.Min(x, 1)
		gc.SetStrokeColor(color.Gray{uint8(x*254 + 1)})
		theta := (0.5 + float64(k)/float64(Orientations)) * math.Pi
		drawOrientedLine(gc, u, v, theta, r)
	}
}

func drawOrientedLine(gc *draw2d.ImageGraphicContext, x, y float64, theta float64, r float64) {
	c := math.Cos(theta)
	s := math.Sin(theta)
	gc.MoveTo(x-r*c, y-r*s)
	gc.LineTo(x+r*c, y+r*s)
	gc.Stroke()
}

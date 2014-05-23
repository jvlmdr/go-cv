package featpyr

import (
	"image"
	"math"
)

type vector struct {
	X, Y float64
}

func vec(p image.Point) vector {
	return vector{float64(p.X), float64(p.Y)}
}

func (p vector) Round() image.Point {
	return image.Pt(round(p.X), round(p.Y))
}

func (p vector) Floor() image.Point {
	return image.Pt(int(math.Floor(p.X)), int(math.Floor(p.Y)))
}

func (p vector) Ceil() image.Point {
	return image.Pt(int(math.Ceil(p.X)), int(math.Ceil(p.Y)))
}

func (p vector) Mul(k float64) vector {
	return vector{k * p.X, k * p.Y}
}

func (p vector) Div(k float64) vector {
	return vector{p.X / k, p.Y / k}
}

func (p vector) Add(q vector) vector {
	return vector{p.X + q.X, p.Y + q.Y}
}

func (p vector) Sub(q vector) vector {
	return vector{p.X - q.X, p.Y - q.Y}
}

func scaleRect(k float64, r image.Rectangle) image.Rectangle {
	a := vec(r.Min).Mul(k).Round()
	b := vec(r.Max).Mul(k).Round()
	return image.Rectangle{a, b}
}

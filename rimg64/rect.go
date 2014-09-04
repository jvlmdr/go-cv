package rimg64

import "image"

// Rect defines an image on an arbitrary rectangular domain.
type Rect struct {
	Image
	Min image.Point
}

func NewRectImage(r image.Rectangle) *Rect {
	im := New(r.Dx(), r.Dy())
	return &Rect{*im, r.Min}
}

func (im *Rect) Bounds() image.Rectangle {
	return image.Rect(0, 0, im.Width, im.Height).Add(im.Min)
}

func (im *Rect) At(i, j int) float64 {
	p := image.Pt(i, j).Sub(im.Min)
	return im.Image.At(p.X, p.Y)
}

func (im *Rect) Set(i, j int, v float64) {
	p := image.Pt(i, j).Sub(im.Min)
	im.Image.Set(p.X, p.Y, v)
}

// RectMulti defines a multi-channel image on an arbitrary rectangular domain.
type RectMulti struct {
	Multi
	Min image.Point
}

func NewRectMulti(r image.Rectangle, channels int) *RectMulti {
	im := NewMulti(r.Dx(), r.Dy(), channels)
	return &RectMulti{*im, r.Min}
}

func (im *RectMulti) Bounds() image.Rectangle {
	return image.Rect(0, 0, im.Width, im.Height).Add(im.Min)
}

func (im *RectMulti) At(i, j, k int) float64 {
	p := image.Pt(i, j).Sub(im.Min)
	return im.Multi.At(p.X, p.Y, k)
}

func (im *RectMulti) Set(i, j, k int, v float64) {
	p := image.Pt(i, j).Sub(im.Min)
	im.Multi.Set(p.X, p.Y, k, v)
}

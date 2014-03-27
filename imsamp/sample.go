package imsamp

import (
	"image"
	"image/color"
	"image/draw"
)

type At func(image.Image, image.Point) color.Color

func Rect(src image.Image, r image.Rectangle, at At) image.Image {
	dst := image.NewRGBA64(image.Rectangle{image.ZP, r.Size()})
	b := dst.Bounds()
	for x := r.Min.X; x < r.Max.X; x++ {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			p := image.Pt(x, y).Sub(r.Min).Add(b.Min)
			dst.Set(p.X, p.Y, at(src, image.Pt(x, y)))
		}
	}
	return dst
}

func Draw(dst image.Image, r image.Rectangle, src image.Image, sp image.Point, at At) image.Image {
	out := dst.(draw.Image)
	for x := r.Min.X; x < r.Max.X; x++ {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			// Remove offset into rectangle and add offset into src.
			p := image.Pt(x, y).Sub(r.Min).Add(sp)
			out.Set(x, y, at(src, p))
		}
	}
	return dst
}

func Continue(im image.Image, p image.Point) color.Color {
	b := im.Bounds()
	x := clip(p.X, b.Min.X, b.Max.X)
	y := clip(p.Y, b.Min.Y, b.Max.Y)
	return im.At(x, y)
}

// Clips x to be in [a, b).
func clip(x, a, b int) int {
	if x < a {
		return a
	}
	if x > b-1 {
		return b - 1
	}
	return x
}

func Black(im image.Image, p image.Point) color.Color {
	if p.In(im.Bounds()) {
		return im.At(p.X, p.Y)
	}
	return color.Black
}

func White(im image.Image, p image.Point) color.Color {
	if p.In(im.Bounds()) {
		return im.At(p.X, p.Y)
	}
	return color.White
}

func Periodic(im image.Image, p image.Point) color.Color {
	p = p.Mod(im.Bounds())
	return im.At(p.X, p.Y)
}

// Uses half-sample symmetry.
func Symmetric(im image.Image, p image.Point) color.Color {
	b := im.Bounds()
	// Make twice as big.
	d := image.Rectangle{b.Min, b.Max.Add(b.Size())}
	p = p.Mod(d)

	// Move to origin.
	p = p.Sub(b.Min)
	w, h := b.Dx(), b.Dy()
	if p.X > w-1 {
		p.X = 2*w - 1 - p.X
	}
	if p.Y > h-1 {
		p.Y = 2*h - 1 - p.Y
	}
	p = p.Add(b.Min)

	return im.At(p.X, p.Y)
}

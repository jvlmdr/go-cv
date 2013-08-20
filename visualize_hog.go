package cv

import (
	"code.google.com/p/draw2d/draw2d"
	"github.com/jackvalmadre/lin-go/vec"
	"image"
	"image/color"
	"math"
)

func HOGImage(hog RealVectorImage, cellSize int) image.Image {
	width, height := hog.Width*cellSize, hog.Height*cellSize
	pic := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with black.
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			pic.Set(i, j, color.Gray{0})
		}
	}

	// Maximum value in vectorized image.
	max, _ := vec.Max(RealVectorImageAsVector{hog})

	gc := draw2d.NewGraphicContext(pic)
	gc.SetLineWidth(1)

	// Draw cells.
	for x := 0; x < hog.Width; x++ {
		for y := 0; y < hog.Height; y++ {
			drawHOGCell(hog, x, y, gc, cellSize, 0, max)
		}
	}

	return pic
}

func SignedHOGImage(hog RealVectorImage, cellSize int) image.Image {
	width, height := hog.Width*cellSize, hog.Height*cellSize
	pic := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with gray.
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			pic.Set(i, j, color.Gray{128})
		}
	}

	// Maximum absolute value in vectorized image.
	max := vec.InfNorm(RealVectorImageAsVector{hog})

	gc := draw2d.NewGraphicContext(pic)
	gc.SetLineWidth(2)

	// Draw cells.
	for x := 0; x < hog.Width; x++ {
		for y := 0; y < hog.Height; y++ {
			drawHOGCell(hog, x, y, gc, cellSize, -max, max)
		}
	}

	return pic
}

func drawHOGCell(hog RealVectorImage, i, j int, gc *draw2d.ImageGraphicContext, cellSize int, min, max float64) {
	u := (float64(i) + 0.5) * float64(cellSize)
	v := (float64(j) + 0.5) * float64(cellSize)
	r := float64(cellSize) / 2

	offset := 0
	if hog.Channels == 31 {
		offset += 18
	}

	for k := 0; k < HOGOrientations; k++ {
		x := (hog.At(i, j, offset+k) - min) / (max - min)
		x = math.Max(x, 0)
		x = math.Min(x, 1)
		gc.SetStrokeColor(color.Gray{uint8(x*254 + 1)})
		theta := (0.5 + float64(k)/float64(HOGOrientations)) * math.Pi
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

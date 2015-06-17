package detect

import (
	"image"
	"math"

	"github.com/jvlmdr/go-cv/feat"
)

type PadRect struct {
	// Size of region.
	Size image.Point
	// Interior of region.
	Int image.Rectangle
}

// Creates a padded rectangle whose interior is size with the given margin.
func Pad(size image.Point, m feat.Margin) PadRect {
	return PadRect{
		Size: image.Pt(size.X+m.Left+m.Right, size.Y+m.Top+m.Bottom),
		Int:  image.Rectangle{Max: size}.Add(m.TopLeft()),
	}
}

func (p PadRect) Left() int   { return p.Int.Min.X }
func (p PadRect) Right() int  { return p.Size.X - p.Int.Max.X }
func (p PadRect) Top() int    { return p.Int.Min.Y }
func (p PadRect) Bottom() int { return p.Size.Y - p.Int.Max.Y }

// Takes a bounding box in an image r.
// Coerces it to the aspect ratio of target.Int according to mode.
// Returns the rectangle which, when resized to target.Size,
// will have the bounding box in target.Int.
// Panics if the target has non-positive interior.
func FitRect(orig image.Rectangle, target PadRect, mode string) (scale float64, fit image.Rectangle) {
	if target.Int.Dx() <= 0 || target.Int.Dy() <= 0 {
		panic("empty interior")
	}
	if mode == "stretch" {
		panic("stretch mode not supported by FitRect")
	}
	aspect := float64(target.Int.Dx()) / float64(target.Int.Dy())
	// Width and height of box in image.
	w, h := float64(orig.Dx()), float64(orig.Dy())
	// Co-erce size to match aspect ratio.
	w, h = SetAspect(w, h, aspect, mode)
	// If source is smaller than target, then scale is > 1 (i.e. need to magnify).
	scale = float64(target.Int.Dx()) / w // == float64(target.Int.Dy()) / h
	// Get position of interior centroid in target rectangle.
	left, top := centroid(target.Int)
	right, bottom := float64(target.Size.X)-left, float64(target.Size.Y)-top
	// Get position of centroid of original bounding box in image.
	x, y := centroid(orig)
	// Scale offsets on all sides and add to centroid for final rectangle.
	// If scale is greater than 1 then source is smaller than target.
	// Then the rectangle in the source image is shrunk (i.e. divide by scale).
	x0, x1 := x-left/scale, x+right/scale
	y0, y1 := y-top/scale, y+bottom/scale
	fit = image.Rect(round(x0), round(y0), round(x1), round(y1))
	return
}

// Change the aspect ratio of a rectangle.
// The mode can be "area", "width", "height", "fit", "fill" or "stretch".
// Panics if mode is empty or unrecognized.
func SetAspect(w, h, aspect float64, mode string) (float64, float64) {
	switch mode {
	case "area":
		// aspect = width / height
		// width = height * aspect
		// width^2 = width * height * aspect
		// height = width / aspect
		// height^2 = width * height / aspect
		w, h = math.Sqrt(w*h*aspect), math.Sqrt(w*h/aspect)
	case "width":
		// Set height from width.
		h = w / aspect
	case "height":
		// Set width from height.
		w = h * aspect
	case "fit":
		// Shrink one dimension.
		w, h = math.Min(w, h*aspect), math.Min(h, w/aspect)
	case "fill":
		// Grow one dimension.
		w, h = math.Max(w, h*aspect), math.Max(h, w/aspect)
	case "stretch":
		// Do nothing.
	case "":
		panic("no mode specified")
	default:
		panic("unknown mode: " + mode)
	}
	return w, h
}

package cv

import (
	"github.com/jackvalmadre/go-vec"
	"image"
	"image/color"
)

// Describes an image with real scalar values.
type RealImage struct {
	Pixels []float64
	Width  int
	Height int
}

func NewRealImage(width, height int) RealImage {
	pixels := make([]float64, width*height)
	return RealImage{pixels, width, height}
}

func (f RealImage) Empty() bool {
	return len(f.Pixels) == 0
}

func (f RealImage) Size() image.Point {
	return image.Pt(f.Width, f.Height)
}

func (f RealImage) At(x, y int) float64 {
	// "Basis vectors" for addressing pixels.
	i := f.Height
	j := 1
	return f.Pixels[x*i+y*j]
}

func (f RealImage) Set(x, y int, v float64) {
	i := f.Height
	j := 1
	f.Pixels[x*i+y*j] = v
}

func (f RealImage) Copy() RealImage {
	g := NewRealImage(f.Width, f.Height)
	copy(g.Pixels, f.Pixels)
	return g
}

func RealImageToGray(f RealImage) image.Gray {
	im := *image.NewGray(image.Rect(0, 0, f.Width, f.Height))

	for x := 0; x < f.Width; x += 1 {
		for y := 0; y < f.Height; y += 1 {
			im.SetGray(x, y, color.Gray{uint8(f.At(x, y) * 255)})
		}
	}

	return im
}

// Describes an image with real vector values.
type RealVectorImage struct {
	Pixels   []float64
	Width    int
	Height   int
	Channels int
}

func NewRealVectorImage(width, height, channels int) RealVectorImage {
	pixels := make([]float64, width*height*channels)
	return RealVectorImage{pixels, width, height, channels}
}

func (f RealVectorImage) Empty() bool {
	return len(f.Pixels) == 0
}

func (f RealVectorImage) ImageSize() image.Point {
	return image.Pt(f.Width, f.Height)
}

func (f RealVectorImage) At(x, y, d int) float64 {
	i := f.Channels * f.Height
	j := f.Channels
	k := 1
	return f.Pixels[x*i+y*j+d*k]
}

func (f RealVectorImage) Set(x, y, d int, v float64) {
	i := f.Channels * f.Height
	j := f.Channels
	k := 1
	f.Pixels[x*i+y*j+d*k] = v
}

func (f RealVectorImage) Copy() RealVectorImage {
	g := NewRealVectorImage(f.Width, f.Height, f.Channels)
	copy(g.Pixels, f.Pixels)
	return g
}

func (src RealVectorImage) CopyChannels(channels []int) RealVectorImage {
	dst := NewRealVectorImage(src.Width, src.Height, len(channels))
	for i := 0; i < src.Width; i++ {
		for j := 0; j < src.Height; j++ {
			for k, p := range channels {
				dst.Set(i, j, k, src.At(i, j, p))
			}
		}
	}
	return dst
}

func (f RealVectorImage) NormalizePositive() {
	x := RealVectorImageAsVector{f}
	max := vec.Max(x)
	vec.ScaleAndCopyTo(x, 1/max, x)
}

func ColorImageToReal(im image.Image) RealVectorImage {
	width, height := im.Bounds().Dx(), im.Bounds().Dy()
	f := NewRealVectorImage(width, height, 3)

	for x := 0; x < width; x++ {
		u := x + im.Bounds().Min.X
		for y := 0; y < height; y++ {
			v := y + im.Bounds().Min.Y

			var c [3]uint32
			c[0], c[1], c[2], _ = im.At(u, v).RGBA()

			for d := 0; d < 3; d++ {
				f.Set(x, y, d, float64(c[d])/float64(0xFFFF))
			}
		}
	}

	return f
}

// Accesses one dimension of a vector-valued image as a scalar image.
type SliceOfRealVectorImage struct {
	Image   RealVectorImage
	Channel int
}

func (slice SliceOfRealVectorImage) Size() image.Point {
	return image.Pt(slice.Image.Width, slice.Image.Height)
}

func (slice SliceOfRealVectorImage) At(x, y int) float64 {
	return slice.Image.At(x, y, slice.Channel)
}

func (slice SliceOfRealVectorImage) Set(x, y int, v float64) {
	slice.Image.Set(x, y, slice.Channel, v)
}

func (f RealVectorImage) Channel(d int) SliceOfRealVectorImage {
	return SliceOfRealVectorImage{f, d}
}

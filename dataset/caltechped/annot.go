package caltechped

import (
	"errors"
	"fmt"
	"image"
	"math"
	"unicode/utf8"

	"github.com/jvlmdr/go-file/fileutil"
)

// ImageAnnot is the annotation of an image.
type ImageAnnot struct {
	Objects []Object
}

// Object is an object within an image.
type Object struct {
	Label string
	Rect  image.Rectangle
	// If Occl is true then Vis contains the visible part.
	Occl bool
	Vis  image.Rectangle
}

func (obj Object) VisFrac() float64 {
	if !obj.Occl {
		return 1
	}
	frac := float64(area(obj.Vis)) / float64(area(obj.Rect))
	return math.Min(1, frac)
}

// ObjectFilter decides whether to use a ground truth annotation.
type ObjectFilter func(obj Object) bool

func Reasonable(obj Object) bool {
	if obj.Label != "person" {
		return false
	}
	// Minimum height 50 pixels.
	if obj.Rect.Dy() < 50 {
		return false
	}
	// At least 65% visible.
	if obj.VisFrac() < 0.65 {
		return false
	}
	return true
}

// LoadAnnot loads the annotation of an image.
// It discards labels that do not satisfy KeepLabel.
func LoadAnnot(fname string) (ImageAnnot, error) {
	lines, err := fileutil.LoadLines(fname)
	if err != nil {
		return ImageAnnot{}, err
	}
	var objs []Object
	for _, line := range lines {
		if len(line) == 0 {
			continue // Empty line.
		}
		if firstRune(line) == '%' {
			continue // Comment.
		}
		raw, err := parseObject(line)
		if err != nil {
			return ImageAnnot{}, fmt.Errorf("parse line: %v", err)
		}
		obj, err := raw.Object()
		if err != nil {
			return ImageAnnot{}, err
		}
		if !KeepLabel(obj.Label) {
			continue
		}
		objs = append(objs, obj)
	}
	return ImageAnnot{Objects: objs}, nil
}

func firstRune(s string) rune {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size == 1 {
		panic("invalid string")
	}
	return r
}

// KeepLabel determines whether to include a box in the image annotation.
func KeepLabel(label string) bool {
	// TODO: Sure not to include "person-fa" in this set?
	switch label {
	case "person", "people", "person?", "ignore":
		return true
	default:
		return false
	}
}

type object struct {
	lbl string
	bb  rect
	occ int
	vbb rect
	ign int
	ang int
}

func (obj object) Object() (Object, error) {
	if obj.ang != 0 {
		return Object{}, fmt.Errorf("angle is not zero: %d", obj.ang)
	}
	if obj.ign != 0 {
		return Object{}, errors.New("object is ignored")
	}
	return Object{obj.lbl, obj.bb.Rect(), obj.occ != 0, obj.vbb.Rect()}, nil
}

type rect struct {
	Left, Top, Width, Height int
}

func (r rect) Rect() image.Rectangle {
	return image.Rect(r.Left, r.Top, r.Left+r.Width, r.Top+r.Height)
}

func parseObject(line string) (object, error) {
	var x object
	n, err := fmt.Sscanf(line, "%s %d %d %d %d %d %d %d %d %d %d %d",
		&x.lbl,
		&x.bb.Left, &x.bb.Top, &x.bb.Width, &x.bb.Height,
		&x.occ,
		&x.vbb.Left, &x.vbb.Top, &x.vbb.Width, &x.vbb.Height,
		&x.ign, &x.ang,
	)
	if err != nil {
		return object{}, err
	}
	if n != 12 {
		return object{}, fmt.Errorf("scan %d arguments, expect 12", n)
	}
	return x, nil
}

func area(rect image.Rectangle) int {
	return rect.Dx() * rect.Dy()
}

func aspect(rect image.Rectangle) float64 {
	return float64(rect.Dx()) / float64(rect.Dy())
}

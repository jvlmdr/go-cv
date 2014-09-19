package main

import (
	"errors"
	"fmt"
	"image"
	"log"
	"math"
	"unicode/utf8"

	"github.com/jvlmdr/go-file/fileutil"
)

type Annot struct {
	Objects []Object
}

type Object struct {
	Label string
	Rect  image.Rectangle
	Occl  bool
	Vis   image.Rectangle
}

func (obj Object) Visible() float64 {
	if !obj.Occl {
		return 1
	}
	frac := float64(area(obj.Vis)) / float64(area(obj.Rect))
	if frac > 1 {
		const eps = 1e-9
		if frac <= 1+eps {
			return 1
		}
		panic(fmt.Sprintf("visible fraction cannot exceed one: %g", frac))
	}
	return frac
}

func Reasonable() *AnnotFilter {
	return &AnnotFilter{
		Labels:  []string{"person"},
		Ignore:  []string{"people", "person?", "person-fa"},
		Height:  Interval{Min: 50, Max: math.Inf(1)},
		Visible: Interval{Min: 0.65, Max: 1},
		Aspect:  Interval{Min: 0, Max: math.Inf(1)},
	}
}

type AnnotFilter struct {
	Labels  StrSet
	Ignore  StrSet
	Height  Interval
	Visible Interval
	Aspect  Interval
}

type StrSet []string

func (set StrSet) Contains(x string) bool {
	for _, s := range set {
		if x == s {
			return true
		}
	}
	return false
}

// Defines an interval [Min, Max] or its complement.
type Interval struct {
	Min float64
	Max float64
	// Take the complement of the interval.
	Inv bool
}

func (r Interval) Contains(x float64) bool {
	in := r.Min <= x && x <= r.Max
	if r.Inv {
		in = !in
	}
	return in
}

// Rects returns the reference and ignore rectangles for the "person" class.
func Rects(a *Annot, filt *AnnotFilter) (refs, ignore []image.Rectangle) {
	for _, obj := range a.Objects {
		switch filt.apply(obj) {
		case decKeep:
			refs = append(refs, obj.Rect)
		case decIgnore:
			ignore = append(ignore, obj.Rect)
		}
	}
	return refs, ignore
}

type filterDecision int

const (
	decKeep filterDecision = iota
	decDiscard
	decIgnore
)

func (filt *AnnotFilter) apply(obj Object) filterDecision {
	if !filt.Labels.Contains(obj.Label) {
		// Object is not in positive class.
		if filt.Ignore.Contains(obj.Label) {
			// Object is in ignored class.
			return decIgnore
		}
		return decDiscard
	}
	// Object is in positive class.
	if !filt.Visible.Contains(obj.Visible()) {
		// Object has visibility outside range.
		return decIgnore
	}
	if !filt.Height.Contains(float64(obj.Rect.Dy())) {
		// Object has height outside range.
		return decIgnore
	}
	if !filt.Aspect.Contains(aspect(obj.Rect)) {
		// Object has height outside range.
		return decIgnore
	}
	if !filt.Aspect.Contains(aspect(obj.Rect)) {
		// Object has height outside range.
		return decIgnore
	}
	return decKeep
}

func loadAnnot(fname string) (*Annot, error) {
	lines, err := fileutil.LoadLines(fname)
	if err != nil {
		return nil, err
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
			return nil, fmt.Errorf("parse line: %v", err)
		}
		obj, err := raw.Object()
		if err != nil {
			return nil, err
		}
		objs = append(objs, obj)
	}
	return &Annot{Objects: objs}, nil
}

func firstRune(s string) rune {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size == 1 {
		panic("invalid string")
	}
	return r
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
		log.Println("error: parse:", line)
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

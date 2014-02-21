package featpyr

import (
	"github.com/jackvalmadre/go-cv/imgpyr"
	"github.com/jackvalmadre/go-cv/rimg64"

	"container/list"
	"image"
	"log"
	"sort"
)

// Need pyramid for scales and rate.
func suppressTable(resps []*rimg64.Image, cands []pyrdet, pyr *Pyramid, pixsize image.Point, maxnum int, maxinter float64) []pyrdet {
	// Sort this list by score.
	sort.Sort(sort.Reverse(byScore(cands)))
	log.Print("sorted candidates by score")

	rem := newLookup(resps, cands, pyr, pixsize)

	// Select best detection, remove those which overlap with it.
	var dets []pyrdet
	for rem.List.Len() > 0 && (maxnum <= 0 || len(dets) < maxnum) {
		dets = append(dets, rem.Pop(maxinter))
	}
	return dets
}

type lookup struct {
	List    *list.List
	Table   [][][]*list.Element
	Pyr     *Pyramid
	PixSize image.Point
}

func newLookup(resps []*rimg64.Image, dets []pyrdet, pyr *Pyramid, pixsize image.Point) *lookup {
	// Linked-list of remaining detections, ordered by score.
	rem := list.New()
	// Look-up of list elements by position.
	table := make([][][]*list.Element, len(resps))
	for k, resp := range resps {
		table[k] = make([][]*list.Element, resp.Width)
		for x := range table[k] {
			table[k][x] = make([]*list.Element, resp.Height)
		}
	}
	// Populate both.
	for _, det := range dets {
		table[det.Level][det.Pos.X][det.Pos.Y] = rem.PushBack(det)
	}
	return &lookup{rem, table, pyr, pixsize}
}

func (rem *lookup) Contains(k, x, y int) bool {
	return rem.Table[k][x][y] != nil
}

func (rem *lookup) Remove(k, x, y int) {
	if rem.Table[k][x][y] == nil {
		panic("remove: element not present")
	}
	rem.List.Remove(rem.Table[k][x][y])
	rem.Table[k][x][y] = nil
}

func (rem *lookup) Pop(maxinter float64) pyrdet {
	// Remove from remaining and add to detections.
	e := rem.List.Front()
	det := e.Value.(pyrdet)
	// Remove.
	rem.Remove(det.Level, det.Pos.X, det.Pos.Y)

	// Scale-space position in feature pyramid.
	p := det.Point
	// 2D position in pixels at level 0.
	c0 := vec(p.Pos.Mul(rem.Pyr.Rate)).Mul(1 / rem.Pyr.Scale(p.Level))

	for k := range rem.Table {
		// 2D position in pixels at level k.
		ck := c0.Mul(rem.Pyr.Scale(k))
		// Domain of response image.
		resp := image.Rect(0, 0, len(rem.Table[k]), len(rem.Table[k][0]))
		// Bounds for search. 2D positions in feature image at level k.
		a := ck.Sub(vec(rem.PixSize)).Div(float64(rem.Pyr.Rate)).Floor()
		b := ck.Add(vec(rem.PixSize)).Div(float64(rem.Pyr.Rate)).Ceil()
		bnds := image.Rectangle{a, b}
		// Restrict to region inside response image.
		bnds = bnds.Intersect(resp)

		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
				if !rem.Contains(k, x, y) {
					// Element already removed.
					continue
				}
				// Current scale-space position in feature pyramid.
				q := imgpyr.Point{k, image.Pt(x, y)}
				rp := rem.Pyr.rectAt(p, rem.PixSize)
				rq := rem.Pyr.rectAt(q, rem.PixSize)
				if intersect(rp, rq, maxinter) {
					// Remove.
					rem.Remove(k, x, y)
				}
			}
		}
	}
	return det
}

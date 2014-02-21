package featpyr

import (
	"github.com/jackvalmadre/go-cv/rimg64"

	"container/list"
	"image"
	"log"
	"sort"
)

func suppressExhaust(resps []*rimg64.Image, cands []pyrdet, pyr *Pyramid, pixsize image.Point, maxnum int, maxinter float64) []pyrdet {
	// Sort by score.
	sort.Sort(sort.Reverse(byScore(cands)))
	log.Print("sorted candidates by score")
	// Copy into linked list.
	rem := list.New()
	for _, det := range cands {
		rem.PushBack(det)
	}
	// Select best detection, remove those which overlap with it.
	var dets []pyrdet
	for rem.Len() > 0 && (maxnum <= 0 || len(dets) < maxnum) {
		det := popList(rem, pyr, pixsize, maxinter)
		dets = append(dets, det)
	}
	return dets
}

func popList(rem *list.List, pyr *Pyramid, pixsize image.Point, maxinter float64) pyrdet {
	det := rem.Remove(rem.Front()).(pyrdet)
	var next *list.Element
	for e := rem.Front(); e != nil; e = next {
		// Buffer next so that we can remove e.
		next = e.Next()
		// Get candidate detection.
		cand := e.Value.(pyrdet)
		rp := pyr.rectAt(det.Point, pixsize)
		rq := pyr.rectAt(cand.Point, pixsize)
		if intersect(rp, rq, maxinter) {
			// Remove.
			rem.Remove(e)
		}
	}
	return det
}

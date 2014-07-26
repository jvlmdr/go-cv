package main

import (
	"github.com/jvlmdr/go-cv"
	"github.com/jvlmdr/go-cv/imgpyr"

	"container/list"
	"image"
	"math"
	"sort"
)

func nonMaxSupp(scoreImgs []cv.RealImage, scales imgpyr.GeoSeq, size image.Point, maxRelOverlap float64) []PyrPos {
	// Put every window into a list.
	all := allDetections(scoreImgs, scales)
	// Sort this list by score.
	sort.Sort(sort.Reverse(byScore{scoreImgs, all}))

	// List of remaining detections.
	remain := list.New()
	// Look-up of list elements by position.
	elems := make([][][]*list.Element, len(scoreImgs))
	for k, img := range scoreImgs {
		elems[k] = make([][]*list.Element, img.Width)
		for x := 0; x < img.Width; x++ {
			elems[k][x] = make([]*list.Element, img.Height)
		}
	}
	// Populate both.
	for _, det := range all {
		elems[det.Level][det.Pos.X][det.Pos.Y] = remain.PushBack(det)
	}

	// Select best detection, remove those which overlap with it.
	var dets []PyrPos
	for remain.Len() > 0 {
		// Remove from remaining and add to detections.
		e := remain.Front()
		det, ok := e.Value.(PyrPos)
		if !ok {
			panic("unexpected type in list")
		}
		remain.Remove(e)
		dets = append(dets, det)

		// Get bounds of detection in its scale.
		r := image.Rectangle{det.Pos, det.Pos.Add(size)}
		// Scale back into frame of original (feature) image.
		r = scaleRect(1/scales.At(det.Level), r)

		for k, img := range scoreImgs {
			// Scale top-left corner into current level.
			p := scalePoint(scales.At(k), r.Min)
			// Calculate bounds on range (exclusive).
			a := p.Sub(size)
			b := p.Add(size)
			ax, bx := max(a.X+1, 0), min(b.X, img.Width)
			ay, by := max(a.Y+1, 0), min(b.Y, img.Height)

			for x := ax; x < bx; x++ {
				for y := ay; y < by; y++ {
					if elems[k][x][y] == nil {
						// Element already removed.
						continue
					}

					// Candidate rectangle in its scale.
					p := image.Pt(x, y)
					q := image.Rectangle{p, p.Add(size)}
					// Scale back into frame of original (feature) image.
					q = scaleRect(1/scales.At(k), q)

					// Compute overlap.
					overlap := q.Intersect(r)
					relR := float64(area(overlap.Size())) / float64(area(r.Size()))
					relQ := float64(area(overlap.Size())) / float64(area(q.Size()))
					if relR <= maxRelOverlap && relQ <= maxRelOverlap {
						// Maximum relative overlap is not exceeded.
						continue
					}
					// Remove from the list.
					remain.Remove(elems[k][x][y])
					elems[k][x][y] = nil
				}
			}
		}
	}
	return dets
}

func scalePoint(k float64, p image.Point) image.Point {
	x := round(k * float64(p.X))
	y := round(k * float64(p.Y))
	return image.Pt(x, y)
}

func scaleRect(k float64, r image.Rectangle) image.Rectangle {
	return image.Rectangle{scalePoint(k, r.Min), scalePoint(k, r.Max)}
}

func area(p image.Point) int {
	return p.X * p.Y
}

func allDetections(scoreImgs []cv.RealImage, scales imgpyr.GeoSeq) []PyrPos {
	// Put every window into a list and sort by score.
	var dets []PyrPos
	for k, scores := range scoreImgs {
		for x := 0; x < scores.Width; x++ {
			for y := 0; y < scores.Height; y++ {
				det := PyrPos{k, image.Pt(x, y)}
				dets = append(dets, det)
			}
		}
	}
	return dets
}

type byScore struct {
	Scores []cv.RealImage
	Dets   []PyrPos
}

func (s byScore) Len() int { return len(s.Dets) }

func (s byScore) Less(i, j int) bool {
	p := s.Dets[i]
	q := s.Dets[j]
	a := s.Scores[p.Level].At(p.Pos.X, p.Pos.Y)
	b := s.Scores[q.Level].At(q.Pos.X, q.Pos.Y)
	return a < b
}

func (s byScore) Swap(i, j int) {
	s.Dets[i], s.Dets[j] = s.Dets[j], s.Dets[i]
}

func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

func round(x float64) int {
	return int(math.Floor(x + 0.5))
}

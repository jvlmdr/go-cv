package inria

import (
	"fmt"
	"image"
	"log"
)

// Annotations of people within an image.
type Annot struct {
	// Image path relative to INRIA root.
	Image string
	// Bounding boxes of pedestrians in image.
	Rects []image.Rectangle
}

// Load image annotation from file.
func LoadAnnot(fname string) (Annot, error) {
	lines, err := loadLines(fname)
	if err != nil {
		return Annot{}, err
	}
	return parseAnnot(lines), nil
}

// Parse image annotation from lines of a file.
func parseAnnot(lines []string) Annot {
	// Get image file.
	var im string
	for _, line := range lines {
		n, _ := fmt.Sscanf(line, "Image filename : %q", &im)
		if n == 1 {
			break
		}
	}

	// Get object bounds.
	var rects []image.Rectangle
	for _, line := range lines {
		var class string
		var xmin, ymin, xmax, ymax int
		const format = "Bounding box for object %d %q (Xmin, Ymin) - (Xmax, Ymax) : (%d, %d) - (%d, %d)"
		n, _ := fmt.Sscanf(line, format, new(int), &class, &xmin, &ymin, &xmax, &ymax)
		if n != 6 {
			continue
		}
		if class != "PASperson" {
			log.Println("found non-person:", class)
		}
		rects = append(rects, image.Rect(xmin, ymin, xmax, ymax))
	}
	return Annot{im, rects}
}

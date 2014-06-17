package main

import (
	"path"

	"github.com/jackvalmadre/go-cv/dataset/inria"
)

func loadAnnots(fname, dir string) ([]inria.Annot, error) {
	files, err := loadLines(fname)
	if err != nil {
		return nil, err
	}

	annots := make([]inria.Annot, len(files))
	for i := range files {
		annot, err := inria.LoadAnnot(path.Join(dir, files[i]))
		if err != nil {
			return nil, err
		}
		annots[i] = annot
	}
	return annots, nil
}

// Converts a list of images to a list of annotations with no rectangles.
// Used to test the negative set.
func imsToAnnots(ims []string) []inria.Annot {
	annots := make([]inria.Annot, len(ims))
	for i, im := range ims {
		annots[i] = inria.Annot{im, nil}
	}
	return annots
}

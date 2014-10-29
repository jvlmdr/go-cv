package voc

import (
	"encoding/xml"
	"image"
	"os"
	"path"
)

// Annot describes the annotation of an entire image.
type Annot struct {
	Image   string
	Objects []Object
}

// Object describes an instance of a class in an image.
type Object struct {
	Class  string
	Region image.Rectangle
	// Extra options.
	Tags
}

type Tags struct {
	Difficult bool
	Occluded  bool
	Truncated bool
}

// LoadAnnot loads the annotation of one image.
//
// Looks in <dir>/Annotations/<image>.xml.
func LoadAnnot(dir, im string) (Annot, error) {
	// Open file.
	name := path.Join(dir, "Annotations", im+".xml")
	fi, err := os.Open(name)
	if err != nil {
		return Annot{}, err
	}
	defer fi.Close()

	// Parse from XML.
	var data struct {
		XMLName xml.Name `xml:"annotation"`
		Objects []struct {
			Name   string `xml:"name"`
			BndBox struct {
				XMin int `xml:"xmin"`
				YMin int `xml:"ymin"`
				XMax int `xml:"xmax"`
				YMax int `xml:"ymax"`
			} `xml:"bndbox"`
			Difficult bool `xml:"difficult"`
			Occluded  bool `xml:"occluded"`
			Truncated bool `xml:"truncated"`
		} `xml:"object"`
	}
	if err := xml.NewDecoder(fi).Decode(&data); err != nil {
		return Annot{}, err
	}

	// Construct from XML object.
	objs := make([]Object, len(data.Objects))
	for i, raw := range data.Objects {
		box := raw.BndBox
		obj := Object{
			raw.Name,
			image.Rect(box.XMin, box.YMin, box.XMax, box.YMax),
			Tags{raw.Difficult, raw.Occluded, raw.Truncated},
		}
		objs[i] = obj
	}
	imfile := imageFile(im)
	return Annot{Image: imfile, Objects: objs}, nil
}

// Returns <dir>/JPEGImages/<im>.jpg.
func imageFile(im string) string {
	return path.Join("JPEGImages", im+".jpg")
}

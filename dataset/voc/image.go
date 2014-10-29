package voc

import (
	"fmt"
	"path"
	"strconv"
	"strings"
)

// LoadImages loads the list of all images in the given set.
//
// Reads <dir>/ImageSets/Main/<set>.txt.
func LoadImages(dir, set string) ([]string, error) {
	fname := path.Join("ImageSets", "Main", set+".txt")
	return loadLines(path.Join(dir, fname))
}

// LoadPosImages loads the list of all images containing objects of the given class in a set.
//
// Reads <dir>/ImageSets/Main/<class>_<set>.txt.
func LoadPosImages(dir, set, class string) ([]string, error) {
	return imagesSetClassLabel(dir, set, class, 0, 2)
}

// LoadNegImages loads the list of images which are guaranteed not to contain objects of the given class.
//
// Reads <dir>/ImageSets/Main/<class>_<set>.txt.
func LoadNegImages(dir, set, class string) ([]string, error) {
	return imagesSetClassLabel(dir, set, class, -1, 0)
}

// Loads list of all image names.
//
// Reads <dir>/ImageSets/Main/<class>_<set>.txt.
func imagesSetClassLabel(dir, set, class string, a, b int) ([]string, error) {
	fname := path.Join("ImageSets", "Main", class+"_"+set+".txt")
	lines, err := loadLines(path.Join(dir, fname))
	if err != nil {
		return nil, err
	}
	var ims []string
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("line does not have 2 fields: %s", line)
		}
		label64, err := strconv.ParseInt(fields[1], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("parse label: %v", err)
		}
		label := int(label64)
		if err := validateLabel(label); err != nil {
			return nil, err
		}
		if !(a <= label && label < b) {
			continue
		}
		ims = append(ims, fields[0])
	}
	return ims, nil
}

func validateLabel(l int) error {
	if l > 1 || l < -1 {
		return fmt.Errorf("label is not in valid range: %d", l)
	}
	return nil
}

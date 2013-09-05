// HOG implementation of Felzenszwalb, Girshick, McAllester, Ramanan (FGMR).

package hog

// #cgo CFLAGS: -Wall -Werror
// #cgo LDFLAGS: -lm
// #include "hog.h"
import "C"

import (
	"github.com/jackvalmadre/go-cv"

	"unsafe"
)

const Orientations = 9
const Channels = 3*Orientations + 4

func HOG(im cv.RealVectorImage, binSize int) cv.RealVectorImage {
	if im.Channels != 3 {
		panic("Input image must have three channels")
	}
	if binSize < 1 {
		panic("Bin size must be positive")
	}

	// Compute size of HOG image.
	dims := [3]C.int{C.int(im.Height), C.int(im.Width), 3}
	var cells [2]C.int
	var out [3]C.int
	C.size(&dims[0], C.int(binSize), &cells[0], &out[0])

	numCells := cells[0] * cells[1]
	hist := make([]C.double, 18*numCells)
	norm := make([]C.double, numCells)

	// Compute HOG image.
	hog := cv.NewRealVectorImage(int(out[1]), int(out[0]), int(out[2]))
	C.process(
		&dims[0],
		(*C.double)(unsafe.Pointer(&im.Pixels[0])),
		(*C.double)(unsafe.Pointer(&hist[0])),
		(*C.double)(unsafe.Pointer(&norm[0])),
		C.int(binSize),
		&cells[0],
		&out[0],
		(*C.double)(unsafe.Pointer(&hog.Pixels[0])))

	return hog
}

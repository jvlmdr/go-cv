// HOG implementation of Felzenszwalb, Girshick, McAllester, Ramanan (FGMR).

package cv

// #cgo CFLAGS: -Wall -Werror
// #cgo LDFLAGS: -lm
// #include "hog.h"
import "C"

import (
	"unsafe"
)

const HOGOrientations = 9
const HOGChannels = 3*HOGOrientations + 4

func HOG(im RealVectorImage, binSize int) RealVectorImage {
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
	hog := NewRealVectorImage(int(out[1]), int(out[0]), int(out[2]))
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

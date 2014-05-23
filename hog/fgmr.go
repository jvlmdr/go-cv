package hog

// HOG implementation of Felzenszwalb, Girshick, McAllester, Ramanan (FGMR).

// #cgo CFLAGS: -Wall -Werror
// #cgo LDFLAGS: -lm
// #include "fgmr.h"
import "C"

import (
	"unsafe"

	"github.com/jackvalmadre/go-cv/rimg64"
)

const Orientations = 9
const Channels = 3*Orientations + 4

func fgmr(im *rimg64.Multi, sbin int) *rimg64.Multi {
	if im.Channels != 3 {
		panic("Input image must have three channels")
	}
	if sbin < 1 {
		panic("Bin size must be positive")
	}

	// Query size of workspace and output.
	var (
		dims  = [3]C.int{C.int(im.Height), C.int(im.Width), 3}
		cells [2]C.int
		out   [3]C.int
	)
	C.size(&dims[0], C.int(sbin), &cells[0], &out[0])

	var (
		// Allocate output.
		hog = rimg64.NewMulti(int(out[1]), int(out[0]), int(out[2]))
		// Allocate workspace.
		numCells = cells[0] * cells[1]
		hist     = make([]C.double, 18*numCells)
		norm     = make([]C.double, numCells)
	)

	// Compute HOG features.
	C.compute(
		&dims[0],
		(*C.double)(unsafe.Pointer(&im.Elems[0])),
		(*C.double)(unsafe.Pointer(&hist[0])),
		(*C.double)(unsafe.Pointer(&norm[0])),
		C.int(sbin),
		&cells[0],
		&out[0],
		(*C.double)(unsafe.Pointer(&hog.Elems[0])))

	return hog
}

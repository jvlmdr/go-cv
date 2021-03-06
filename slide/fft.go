package slide

import "image"

var primes = []int{2, 3, 5, 7}

func FFTLen(n int) (x, flops int) {
	return fftLen(1, 0, n, primes)
}

func fftLen(prod, sum int, n int, ks []int) (arg, min int) {
	if prod >= n {
		return prod, prod * sum
	}
	for i, k := range ks {
		x, flops := fftLen(prod*k, sum+k, n, ks[i:])
		if i == 0 || flops < min {
			arg, min = x, flops
		}
	}
	return
}

func FFT2Size(n image.Point) (size image.Point, flops int) {
	x, cx := FFTLen(n.X)
	y, cy := FFTLen(n.Y)
	return image.Pt(x, y), x*cy + y*cx
}

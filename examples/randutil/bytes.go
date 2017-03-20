package randutil

import "math/rand"

// usage DataLen: data = data[:randutil.DataLen(r, data)]

func DataLen(r *rand.Rand, data []byte) int {
	x := r.Intn(cap(data))
	shift := 0
	for d := 1; d < x; d <<= 1 {
		shift++
	}
	shift--
	if shift > 0 {
		x >>= uint(r.Intn(shift))
	}
	return x
}

func FillBytes(r *rand.Rand, data []byte) {
	for i := range data {
		data[i] = byte(r.Intn(256))
	}
}

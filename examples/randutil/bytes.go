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

func FillBytes(r *rand.Rand, bs []byte) {
	const (
		bitsPerByte    = 8
		bytesPerUint32 = 4
	)
	var x uint32
	var n int // number of random bytes
	for i := range bs {
		if n == 0 {
			x = r.Uint32()
			n = bytesPerUint32
		}
		bs[i] = byte(x)
		x >>= bitsPerByte
		n--
	}
}

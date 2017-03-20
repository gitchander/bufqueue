package datalen

import (
	"bytes"
	"math/rand"
	"testing"
	"time"
)

func TestLengthSamples(t *testing.T) {

	samples := []struct {
		length int
		data   []byte
	}{
		{0, []byte{0x00}},
		{1, []byte{0x01}},
		{127, []byte{0x7F}},
		{128, []byte{0x80, 0x80}},
		{129, []byte{0x80, 0x81}},
		{8192, []byte{0xA0, 0x00}},
		{16383, []byte{0xBF, 0xFF}},
		{16384, []byte{0xC0, 0x40, 0x00}},
		{16385, []byte{0xC0, 0x40, 0x01}},
		{1048576, []byte{0xD0, 0x00, 0x00}},
		{2097151, []byte{0xDF, 0xFF, 0xFF}},
		{2097152, []byte{0xE0, 0x20, 0x00, 0x00}},
		{2097153, []byte{0xE0, 0x20, 0x00, 0x01}},
		{134217728, []byte{0xE8, 0x00, 0x00, 0x00}},
		{268435455, []byte{0xEF, 0xFF, 0xFF, 0xFF}},
	}

	// encode
	data := make([]byte, MaxSize)
	for _, x := range samples {
		n, err := Encode(x.length, data)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(data[:n], x.data) {
			t.Fatalf("bytes: %x != %x", data[:n], x.data)
		}
	}

	// decode
	var length int
	for _, x := range samples {
		n, err := Decode(x.data, &length)
		if err != nil {
			t.Fatal(err)
		}
		if n != len(x.data) {
			t.Fatalf("lengths: %d != %d", n, len(x.data))
		}
		if length != x.length {
			t.Fatalf("values: %d != %d", length, x.length)
		}
	}
}

func TestLengthRand(t *testing.T) {

	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	randLength := func() int {
		return r.Intn(MaxValue+1) >> uint(r.Intn(7*4))
	}

	for i := 0; i < 1000; i++ {
		var data []byte
		var ls1 = make([]int, r.Intn(50))
		for i := range ls1 {
			length := randLength()
			bs := make([]byte, MaxSize)
			n, err := Encode(length, bs)
			if err != nil {
				t.Fatal(err)
			}
			data = append(data, bs[:n]...)
			ls1[i] = length
		}

		var ls2 []int
		for {
			var length int
			n, err := Decode(data, &length)
			if err != nil {
				if n == 0 {
					break
				}
				t.Fatal(err)
			}
			data = data[n:]
			ls2 = append(ls2, length)
		}
		if len(ls1) != len(ls2) {
			t.Fatalf("len(%d) != len(%d)", len(ls1), len(ls2))
		}
		for i := range ls1 {
			if ls1[i] != ls2[i] {
				t.Fatalf("%d != %d", ls1[i], ls2[i])
			}
		}
	}
}

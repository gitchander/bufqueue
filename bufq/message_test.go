package bufq

import (
	"bytes"
	"io"
	"math/rand"
	"testing"
	"time"
)

func TestMessageRand(t *testing.T) {

	r := randNow()

	as := make([]*Message, r.Intn(1000))
	for i := range as {
		var (
			messageType = randMessageType(r)
			valueLength = r.Intn(65536) >> r.Intn(16)
		)
		data := make([]byte, valueLength)
		randBytes(r, data)

		as[i] = &Message{
			Type:  messageType,
			Value: data,
		}
	}

	var buf bytes.Buffer

	for _, a := range as {
		//t.Logf("message: t=%d data_len=%d\n", a.Type, len(a.Value))
		err := WriteMessage(&buf, a)
		if err != nil {
			t.Fatal(err)
		}
	}

	bs := make([]*Message, 0, len(as))
	for {
		var b = new(Message)
		err := ReadMessage(&buf, b)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
		bs = append(bs, b)
	}

	// t.Logf("as len %d", len(as))
	// t.Logf("bs len %d", len(bs))

	if len(as) != len(bs) {
		t.Fatalf("lengths not equal: %d != %d", len(as), len(bs))
	}

	for i, a := range as {
		b := bs[i]
		if !messagesEqual(a, b) {
			t.Fatalf("messages not equal: %+v != %+v", a, b)
		}
	}
}

func randNow() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func randBool(r *rand.Rand) bool {
	return (r.Int() & 1) == 1
}

func randByte(r *rand.Rand) byte {
	return byte(r.Intn(256))
}

func randBytes(r *rand.Rand, data []byte) {
	for i := range data {
		data[i] = randByte(r)
	}
}

func randMessageType(r *rand.Rand) int {
	x := r.Int() >> r.Intn(64)
	if randBool(r) {
		x = -x
	}
	return x
}

func messagesEqual(a, b *Message) bool {
	if a.Type != b.Type {
		return false
	}
	if !bytes.Equal(a.Value, b.Value) {
		return false
	}
	return true
}

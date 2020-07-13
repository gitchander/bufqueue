package bufq

import (
	"encoding/binary"
	"io"
)

// Message structure (TLV):
// +--------+----------+---------+
// |  Type  |  Length  |  Value  |
// +--------+----------+---------+

type Message struct {
	Type  int
	Value []byte
}

func WriteMessage(w io.Writer, m *Message) error {

	data := make([]byte, binary.MaxVarintLen64)

	// Type
	n := binary.PutVarint(data, int64(m.Type))
	_, err := w.Write(data[:n])
	if err != nil {
		return err
	}

	// Length
	n = binary.PutVarint(data, int64(len(m.Value)))
	_, err = w.Write(data[:n])
	if err != nil {
		return err
	}

	// Value
	_, err = w.Write(m.Value)
	return err
}

func ReadMessage(r io.Reader, m *Message) error {

	br := byteReader{r: r}

	// Type
	x, err := binary.ReadVarint(br)
	if err != nil {
		return err
	}
	messageType := int(x)

	// Length
	x, err = binary.ReadVarint(br)
	if err != nil {
		return err
	}
	length := int(x)

	// Value
	data := m.Value
	if cap(data) < length {
		data = make([]byte, length)
	}
	data = data[:length]
	_, err = io.ReadFull(r, data)
	if err != nil {
		return err
	}

	m.Type = messageType
	m.Value = data

	return nil
}

type byteReader struct {
	data [1]byte
	r    io.Reader
}

var _ io.ByteReader = byteReader{}

func (p byteReader) ReadByte() (byte, error) {
	buf := p.data[:]
	_, err := p.r.Read(buf)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

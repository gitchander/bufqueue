package bufqueue

import (
	"errors"
	"io"
	"time"

	"github.com/gitchander/bufqueue/datalen"
)

type ReaderConfig struct {
	FilesConfig
	ReadDuration time.Duration
}

// Reader is not safe for multiple goroutines!

type Reader struct {
	block []byte
	pos   struct {
		curr int
		prev int
	}

	bfr *blockFileReader
}

func NewReader(rc ReaderConfig) (*Reader, error) {

	bfr, err := newBlockFileReader(rc)
	if err != nil {
		return nil, err
	}

	return &Reader{
		bfr: bfr,
	}, nil
}

func (r *Reader) Close() (err error) {

	if len(r.block) > 0 {
		if r.pos.curr == len(r.block) {
			err = r.bfr.RemoveFile()
			r.prepareBlock(0)
		} else if r.pos.curr > 0 {
			err = r.bfr.Write(r.block[r.pos.curr:])
			r.pos.prev = r.pos.curr
		}
	}

	return
}

func (r *Reader) Read() (message []byte, err error) {
	for {
		message, err = r.readMessage()
		if err == nil {
			return message, nil
		}
		if err == io.EOF {
			fi, err := r.bfr.NextFile()
			if err != nil {
				return nil, err
			}

			r.prepareBlock(int(fi.Size()))
			n, err := r.bfr.Read(r.block)
			if err != nil {
				return nil, err
			}
			if n < len(r.block) {
				return nil, errors.New("short read from file")
			}

			continue
		}
		r.prepareBlock(0)
		r.bfr.MoveFileToBad()
	}
}

func (r *Reader) Unread() error {
	if r.pos.curr > r.pos.prev {
		r.pos.curr = r.pos.prev
		return nil
	}
	return errors.New("bufqueue: last operation is not read")
}

func (r *Reader) readMessage() (message []byte, err error) {

	data := r.block[r.pos.curr:]

	if len(data) == 0 {
		return nil, io.EOF
	}

	var length int

	l_size, err := datalen.Decode(data, &length)
	if err != nil {
		return nil, err
	}

	recLen := l_size + length
	if len(data) < recLen {
		return nil, errors.New("bufqueue: read record: insufficient data length")
	}

	message = make([]byte, length)
	copy(message, data[l_size:])

	r.pos.prev = r.pos.curr
	r.pos.curr += recLen

	return message, nil
}

func (r *Reader) prepareBlock(size int) {

	if cap(r.block) < size {
		r.block = make([]byte, size)
	}
	r.block = r.block[:size]

	r.pos.curr = 0
	r.pos.prev = 0
}

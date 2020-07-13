package bufq

import (
	"bytes"
	"errors"
	"io"
	"sync"
	"time"
)

// io.Pipe

var ErrClosedPipe = errors.New("bufq: read/write on closed pipe")

type BlockConfig struct {
	MaxSize          int // max size of the block
	MessagesPerBlock int
}

type Config struct {
	BlockConfig
	FilesConfig  FilesConfig
	FlushTimeout time.Time
}

type pipe struct {
	messages chan *Message

	once sync.Once // Protects closing done
	done chan struct{}

	rerr syncError
	werr syncError
}

func (p *pipe) CloseRead(err error) error {
	if err == nil {
		err = ErrClosedPipe
	}
	p.rerr.Store(err)
	p.once.Do(func() { close(p.done) })
	return nil
}

func (p *pipe) CloseWrite(err error) error {
	if err == nil {
		err = io.EOF
	}
	p.werr.Store(err)
	p.once.Do(func() { close(p.done) })
	return nil
}

func (p *pipe) WriteMessage(m *Message) error {

	select {
	case <-p.done:
		return ErrClosedPipe

	case p.messages <- m:

	}

	return nil
}

func (p *pipe) ReadMessage(m *Message) error {

	return nil
}

func (p *pipe) UnreadMessage() error {

	return nil
}

func Pipe(c Config) (*MessageReader, *MessageWriter) {

	messages := make(chan *Message)
	done := make(chan struct{})
	fl := newFileList()

	p := &pipe{
		messages: messages,
		done:     done,
	}

	go messageWriter(done, messages, fl)

	return &MessageReader{p}, &MessageWriter{p}
}

func messageWriter(done <-chan struct{}, messages <-chan *Message, fl *fileList) {

	bw := newBlockWriter()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-(done):
			return

		case <-(ticker.C):
			bw.Flush()

		case message, ok := <-messages:
			if !ok {
				return
			}
			err := bw.WriteMessage(message)
			if err != nil {

			}
		}
	}
}

type blockWriter struct {
	buf bytes.Buffer
}

func newBlockWriter() *blockWriter {
	return nil
}

func (p *blockWriter) Flush() error {

	return nil
}

func (p *blockWriter) WriteMessage(message *Message) error {

	err := WriteMessage(&(p.buf), message)
	if err != nil {
		return err
	}

	return nil
}

package bufqueue

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gitchander/bufqueue/datalen"
)

var errWriterClosed = errors.New("bufqueue: writer is closed")

type BlockConfig struct {
	MaxSize         int           // max block size (in bytes)
	MaxRecordsCount int           // max records count in block
	Duration        time.Duration // duration for write block
}

type WriterConfig struct {
	FilesConfig
	Block BlockConfig
}

// Writer is safe for multiple goroutines!

type Writer struct {
	mutex        sync.Mutex
	opened       bool
	event        chan bool
	block        []byte
	lengthData   []byte // use for encode length
	recordsCount int
	wc           WriterConfig
	bfw          *blockFileWriter
}

func NewWriter(wc WriterConfig) (*Writer, error) {

	bfw, err := newBlockFileWriter(wc.FilesConfig)
	if err != nil {
		return nil, err
	}

	w := &Writer{
		opened:     true,
		event:      make(chan bool),
		block:      make([]byte, 0, wc.Block.MaxSize),
		lengthData: make([]byte, datalen.MaxSize),
		wc:         wc,
		bfw:        bfw,
	}

	run := func() {
		w.mutex.Lock()
		defer w.mutex.Unlock()

		err := w.saveBlock()
		if err != nil {
			log.Fatal("bufqueue.Writer: save block error:", err)
		}
	}
	go loopTimeoutRun(w.event, wc.Block.Duration, run)

	return w, nil
}

func (w *Writer) Close() error {

	w.mutex.Lock()
	defer w.mutex.Unlock()

	if !w.opened {
		return errWriterClosed
	}

	w.event <- false // stop timer loop
	<-w.event        // wait done

	err := w.saveBlock()
	w.opened = false
	return err
}

func loopTimeoutRun(event chan bool, dur time.Duration, Run func()) {
	timer := time.NewTimer(dur)
	timer.Stop()
	for {
		select {
		case ok := <-event:
			if ok {
				timer.Reset(dur)
			} else {
				timer.Stop()
				event <- false // done signal
				return
			}
		case <-timer.C:
			Run()
		}
	}
}

func (w *Writer) Write(message []byte) error {

	w.mutex.Lock()
	defer w.mutex.Unlock()

	if !w.opened {
		return errWriterClosed
	}

	l_data := w.lengthData
	l_size, err := datalen.Encode(len(message), l_data)
	if err != nil {
		return err
	}

	recLen := l_size + len(message)

	if recLen > cap(w.block) {
		return errors.New("record size more block size")
	}

	if w.neadSaveBlock(recLen) {
		if err := w.saveBlock(); err != nil {
			return err
		}
	}

	recPos := len(w.block)
	if recPos == 0 {
		w.event <- true // reset timer
	}
	w.block = w.block[:recPos+recLen]

	copy(w.block[recPos:], l_data[:l_size])
	copy(w.block[recPos+l_size:], message)

	w.recordsCount++

	return nil
}

func (w *Writer) neadSaveBlock(recLen int) bool {

	if len(w.block)+recLen > cap(w.block) {
		return true
	}

	if w.recordsCount >= w.wc.Block.MaxRecordsCount {
		return true
	}

	return false
}

func (w *Writer) saveBlock() error {

	if w.recordsCount == 0 {
		return nil
	}

	err := w.bfw.SaveToFile(w.block)

	w.block = w.block[:0]
	w.recordsCount = 0

	return err
}

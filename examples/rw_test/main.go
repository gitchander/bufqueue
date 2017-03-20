package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	bufq "github.com/gitchander/bufqueue"
	"github.com/gitchander/bufqueue/examples/randutil"
)

const (
	kilobyte = 1024
	megabyte = 1024 * kilobyte
)

func main() {

	start := time.Now()

	var (
		hwChan = make(chan []byte)
		hrChan = make(chan []byte)
	)

	fc := bufq.FilesConfig{
		Dirname:    "blocks",
		FilePrefix: "block-",
	}

	go loopWriter(hwChan, 100000, fc)
	go loopReader(hrChan, -1, false, fc)

	hw := <-hwChan
	hr := <-hrChan

	if bytes.Equal(hw, hr) {
		fmt.Println("hashes is equal!")
	} else {
		fmt.Println("hashes is not equal!")
		fmt.Printf("writer hash: %X\n", hw)
		fmt.Printf("reader hash: %X\n", hr)
	}

	fmt.Println(time.Since(start))
}

func loopWriter(hashChan chan<- []byte, n int, fc bufq.FilesConfig) {

	h := sha256.New()
	defer func() {
		hashChan <- h.Sum(nil)
	}()

	wc := bufq.WriterConfig{
		FilesConfig: fc,
		Block: bufq.BlockConfig{
			MaxSize:         10 * megabyte,
			MaxRecordsCount: 10000,
			Duration:        5 * time.Second,
		},
	}

	w, err := bufq.NewWriter(wc)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	r := randutil.NewRandFromTime()
	data := make([]byte, 0, 19000)
	for i := 0; i < n; i++ {
		data = data[:randutil.DataLen(r, data)]
		randutil.FillBytes(r, data)
		if err = w.Write(data); err != nil {
			log.Fatal(err)
		}
		h.Write(data)
	}
}

func loopReader(hashChan chan<- []byte, maxCount int, messageOut bool, fc bufq.FilesConfig) {

	h := sha256.New()
	defer func() {
		hashChan <- h.Sum(nil)
	}()

	rc := bufq.ReaderConfig{
		FilesConfig:  fc,
		ReadDuration: 2 * time.Second,
	}

	r, err := bufq.NewReader(rc)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	count := 0
	for (maxCount == -1) || (count < maxCount) {
		message, err := r.Read()
		if err != nil {
			if err == bufq.ErrReadTimeout {
				break
			}
			log.Fatal(err)
		}
		if messageOut {
			fmt.Println("receive message len:", len(message))
		}
		h.Write(message)
		count++
	}
	fmt.Println("total message count:", count)
}

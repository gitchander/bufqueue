package main

import (
	"flag"
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

	pDirname := flag.String("dir", "test", "dirname")
	pPrefix := flag.String("prefix", "block-", "prefix for block-file")
	pCount := flag.Int("count", 10, "count messages generate")

	flag.Parse()

	n := *pCount

	config := bufq.WriterConfig{
		FilesConfig: bufq.FilesConfig{
			Dirname:    *pDirname,
			FilePrefix: *pPrefix,
		},
		Block: bufq.BlockConfig{
			MaxSize:         10 * megabyte,
			MaxRecordsCount: 10000,
			Duration:        2 * time.Second,
		},
	}

	w, err := bufq.NewWriter(config)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	r := randutil.NewRandNow()
	data := make([]byte, 0, 17000)
	for i := 0; i < n; i++ {
		data = data[:randutil.DataLen(r, data)]
		randutil.FillBytes(r, data)
		if err = w.Write(data); err != nil {
			log.Fatal(err)
		}
	}
}

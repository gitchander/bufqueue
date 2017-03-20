package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	bufq "github.com/gitchander/bufqueue"
)

func main() {

	var (
		pDirname = flag.String("dir", "test", "dirname")
		pPrefix  = flag.String("prefix", "block-", "prefix for block-file")
		pOutlen  = flag.Bool("outlen", false, "print incoming message len")
		pCount   = flag.Int("count", -1, "count messages for read")
	)

	flag.Parse()

	messageLenOut := *pOutlen
	maxCount := *pCount

	config := bufq.ReaderConfig{
		FilesConfig: bufq.FilesConfig{
			Dirname:    *pDirname,
			FilePrefix: *pPrefix,
		},
		ReadDuration: 1 * time.Second,
	}

	r, err := bufq.NewReader(config)
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
		if messageLenOut {
			fmt.Println("message len:", len(message))
		}
		count++
	}
	fmt.Println("total message count:", count)
}

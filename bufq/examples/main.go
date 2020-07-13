package main

import (
	"bytes"
	"fmt"
	"log"

	"bufq"
)

func main() {
	var buf bytes.Buffer

	m := &bufq.Message{
		Type:  0,
		Value: []byte("1"),
	}

	err := bufq.WriteMessage(&buf, m)
	checkError(err)

	fmt.Printf("data hex: [% x]\n", buf.Bytes())

	var m2 = new(bufq.Message)
	err = bufq.ReadMessage(&buf, m2)
	checkError(err)

	fmt.Printf("%+v\n", m2)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

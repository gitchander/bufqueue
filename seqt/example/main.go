package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/gitchander/bufqueue/seqt"
)

func main() {
	simple()
	//exampleParse()
}

func motivation() {
	var a []string
	for i := 0; i < 1000; i++ {
		a = append(a, strconv.Itoa(i))
	}
	sort.Sort(sort.StringSlice(a))
	for _, s := range a {
		fmt.Println(s)
	}
}

func simple() {
	tab := seqt.NewTable(seqt.UPPER_LETTERS)
	seq := new(seqt.Sequence)
	for i := 0; i < 15; i++ {
		seq.Next()
		fmt.Println(tab.String(seq))
	}
}

func exampleParse() {

	tab := seqt.NewTable(seqt.UPPER_LETTERS)

	seq, err := tab.Parse("BBB")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 14; i++ {
		fmt.Println(tab.String(seq))
		seq.Next()
	}
}

func makeStrings() {
	var a []string
	tab := seqt.NewTable(seqt.UPPER_LETTERS)
	seq := new(seqt.Sequence)
	for i := 0; i < 1000; i++ {
		seq.Next()
		s := tab.String(seq)
		a = append(a, s)
	}

	if !Sorted(sort.StringSlice(a)) {
		log.Fatal("sequence not sorted")
	}

	for _, s := range a {
		fmt.Println(s)
	}
}

func makeFiles() {

	tab := seqt.NewTable(seqt.LOWER_LETTERS)
	seq := new(seqt.Sequence)

	data := []byte{0x7A}
	for i := 0; i < 100000; i++ {
		seq.Next()
		s := tab.String(seq)
		err := ioutil.WriteFile(fmt.Sprintf("res/test-%s.log", s), data, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func Sorted(v sort.Interface) bool {
	n := v.Len()
	for i := 1; i < n; i++ {
		if v.Less(i, i-1) {
			return false
		}
	}
	return true
}

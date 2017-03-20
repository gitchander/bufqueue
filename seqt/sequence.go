package seqt

import (
	"errors"

	"github.com/gitchander/minmax"
)

type Sequence struct {
	a []int
}

func (seq *Sequence) Next() {
	seq.a = next(seq.a)
}

func newFromIntSlice(a []int) (*Sequence, error) {

	if len(a) == 0 {
		return new(Sequence), nil
	}

	for i := range a {
		if a[i] < 0 {
			return nil, errors.New("negative value in int slice is not support")
		}
	}

	a = normalIntSlice(a)
	seq := &Sequence{a}

	return seq, nil
}

func normalIntSlice(a []int) []int {

	maxIndex := minmax.IndexOfMax(minmax.IntSlice(a))
	if maxIndex == -1 {
		return a
	}

	maxValue := max(a[maxIndex], len(a)-1)

	for len(a) <= maxValue {
		a = append(a, 0)
	}

	a[maxValue] = maxValue

	return a
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func next(a []int) []int {
	for i := 0; ; i++ {
		if len(a) == i {
			a = append(a, i)
		} else {
			a[i]++
		}
		if a[i] < len(a) {
			break
		}
		a[i] = 0
	}
	return a
}

func value(a []int, table []byte) string {
	n := len(a)
	var data = make([]byte, n)
	for i := range data {
		n--
		data[i] = table[a[n]]
	}
	return string(data)
}

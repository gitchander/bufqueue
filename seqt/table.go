package seqt

import (
	"errors"
	"sort"
)

const (
	DIGITS = 1 << iota
	UPPER_LETTERS
	LOWER_LETTERS
)

type Table struct {
	tag    int
	data   []byte
	shifts struct {
		digits       int
		upperLetters int
		lowerLetters int
	}
}

func NewTable(tag int) *Table {

	t := &Table{tag: tag}

	shift := 0
	var data []byte

	if (tag & DIGITS) != 0 {
		for i := '0'; i <= '9'; i++ {
			data = append(data, byte(i))
		}
		t.shifts.digits = shift
		shift += '9' - '0' + 1
	}

	if (tag & UPPER_LETTERS) != 0 {
		for i := 'A'; i <= 'Z'; i++ {
			data = append(data, byte(i))
		}
		t.shifts.upperLetters = shift
		shift += 'Z' - 'A' + 1
	}

	if (tag & LOWER_LETTERS) != 0 {
		for i := 'a'; i <= 'z'; i++ {
			data = append(data, byte(i))
		}
		t.shifts.lowerLetters = shift
		shift += 'z' - 'a' + 1
	}

	sort.Sort(byteSlice(data))

	t.data = data

	return t
}

var _ sort.IntSlice

type byteSlice []byte

func (x byteSlice) Len() int           { return len(x) }
func (x byteSlice) Less(i, j int) bool { return x[i] < x[j] }
func (x byteSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func (t *Table) parseIntSlice(s string) ([]int, error) {

	bs := []byte(s)

	a := make([]int, len(bs))
	ia := len(a) - 1
	for i, b := range bs {

		var index int

		switch {
		case byteIsDigit(b):
			{
				if (t.tag & DIGITS) == 0 {
					return nil, errors.New("digits is not support")
				}
				index = int(b) - '0' + t.shifts.digits
			}
		case byteIsUpperLetter(b):
			{
				if (t.tag & UPPER_LETTERS) == 0 {
					return nil, errors.New("upper letters is not support")
				}
				index = int(b) - 'A' + t.shifts.upperLetters
			}
		case byteIsLowerLetter(b):
			{
				if (t.tag & LOWER_LETTERS) == 0 {
					return nil, errors.New("lower letters is not support")
				}
				index = int(b) - 'a' + t.shifts.lowerLetters
			}
		default:
			return nil, errors.New("bytes is not support")
		}

		a[ia-i] = index
	}

	return a, nil
}

func (t *Table) Parse(s string) (*Sequence, error) {
	a, err := t.parseIntSlice(s)
	if err != nil {
		return nil, err
	}
	return newFromIntSlice(a)
}

func (t *Table) String(seq *Sequence) string {
	return value(seq.a, t.data)
}

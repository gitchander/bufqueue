package randutil

import (
	"math"
	"math/rand"
)

type alphabet struct {
	lower []rune
	upper []rune
}

var (
	abcEnglish = alphabet{
		lower: []rune("abcdefghijklmnopqrstuvwxyz"),
		upper: []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
	}

	abcRussian = alphabet{
		lower: []rune("абвгдеёжзийклмнопрстуфхчцшщъыьэюя"),
		upper: []rune("АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"),
	}

	abcUkrainian = alphabet{
		lower: []rune("абвгґдеєжзиіїйклмнопрстуфхцчшщьюя"),
		upper: []rune("АБВГҐДЕЄЖЗИІЇЙКЛМНОПРСТУФХЦЧШЩЬЮЯ"),
	}
)

var (
	digits       = []rune("0123456789")
	specialRunes = []rune("!@#$%^&*()_+,./?\"'`")
)

var randRunesSamples = [][]rune{
	digits,
	specialRunes,
	abcEnglish.lower,
	abcEnglish.upper,
	abcRussian.lower,
	abcRussian.upper,
	abcUkrainian.lower,
	abcUkrainian.upper,
}

func RandString(r *rand.Rand, maxLen int) string {

	if maxLen <= 0 {
		return ""
	}

	rs := make([]rune, int(math.Floor(r.Float64()*float64(maxLen))))

	for i := range rs {
		var (
			j = r.Intn(len(randRunesSamples))
			k = r.Intn(len(randRunesSamples[j]))
		)
		rs[i] = randRunesSamples[j][k]
	}

	return string(rs)
}

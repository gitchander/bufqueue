package randutil

import (
	"math/rand"
	"time"
)

func NewRandSeed(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

func NewRandTime(t time.Time) *rand.Rand {
	return NewRandSeed(t.UTC().UnixNano())
}

func NewRandNow() *rand.Rand {
	return NewRandTime(time.Now())
}

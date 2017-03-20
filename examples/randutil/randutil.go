package randutil

import (
	"math/rand"
	"time"
)

func NewRandFromTime() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

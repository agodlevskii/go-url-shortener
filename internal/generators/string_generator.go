package generators

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateString(size int) string {
	if size == 0 {
		size = 7
	}

	r := newRandom()
	b := make([]byte, size)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

func newRandom() *rand.Rand {
	src := rand.NewSource(time.Now().UnixNano())
	return rand.New(src)
}

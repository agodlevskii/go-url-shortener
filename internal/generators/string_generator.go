package generators

import (
	"errors"
	"go-url-shortener/internal/apperrors"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateString(size int) (string, error) {
	if size == 0 {
		return "", errors.New(apperrors.RandomStrLen)
	}

	r := newRandom()
	b := make([]byte, size)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b), nil
}

func newRandom() *rand.Rand {
	src := rand.NewSource(time.Now().UnixNano())
	return rand.New(src)
}

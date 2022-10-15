package generators

import (
	"errors"
	"go-url-shortener/internal/apperrors"
	"math/rand"
	"time"
)

// letterBytes provides a list of the symbols that can be used for the random string generation.
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateString provides a randomly generated string of the required size.
// The generated value only includes the symbols presented in the letterBytes constant.
// The function will return an error if the size is zero.
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

// newRandom provides an element of the rand.Rand type based on the current time.
func newRandom() *rand.Rand {
	src := rand.NewSource(time.Now().UnixNano())
	return rand.New(src)
}

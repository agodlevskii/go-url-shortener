package generators

import (
	"crypto/rand"
	"errors"
	"go-url-shortener/internal/apperrors"
	"log"
	"math/big"
)

// letterBytes provides a list of the symbols that can be used for the random string generation.
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateString provides a randomly generated string of the required size.
// The generated value only includes the symbols presented in the letterBytes constant.
// The function will return an error if the size ilis zero.
func GenerateString(size int) (string, error) {
	if size == 0 {
		return "", errors.New(apperrors.RandomStrLen)
	}

	b := make([]byte, size)
	for i := range b {
		v, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if err != nil {
			log.Fatal(err)
		}

		b[i] = letterBytes[v.Int64()]
	}
	return string(b), nil
}

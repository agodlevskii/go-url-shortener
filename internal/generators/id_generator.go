// Package generators provides the algorithms for the randomly generated values.
package generators

import (
	"errors"
	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/storage"
)

// GenerateID provides a randomly generated ID of the required size.
// The generation algorithm is covered by the GenerateString function.
// The generated ID must be unique for the current DB, presented by the storage.Storager interface.
// The function will return an error if the size is zero.
func GenerateID(db storage.Storager, size int) (string, error) {
	if size == 0 {
		return "", errors.New(apperrors.IDSize)
	}

	for step := 1; step < 10; step++ {
		id, err := GenerateString(size)
		if err != nil {
			return "", err
		}

		has, err := db.Has(id)
		if err != nil {
			return "", err
		}

		if !has {
			return id, nil
		}
	}

	return "", errors.New(apperrors.IDGeneration)
}

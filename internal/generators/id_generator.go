package generators

import (
	"errors"
	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/storage"
)

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

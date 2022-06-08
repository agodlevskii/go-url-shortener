package generators

import (
	"errors"
	"go-url-shortener/internal/storage"
)

func GenerateID(db storage.MemoRepo, size int) (string, error) {
	if size == 0 {
		size = 7
	}

	id := GenerateString(size)

	for step := 1; step < 10; step++ {
		if !db.Has(id) {
			return id, nil
		}
		id = GenerateString(size)
	}

	return "", errors.New("couldn't generate ID")
}

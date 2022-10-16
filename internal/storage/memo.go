package storage

import (
	"errors"
	"go-url-shortener/internal/apperrors"
	"sync"
)

// MemoRepo describes the in-memo implementation of the Storager interface.
// The in-memo storage is implemented via the sync.Map.
type MemoRepo struct {
	db sync.Map
}

// NewMemoryRepo returns a new instance of the MemoRepo type.
func NewMemoryRepo() *MemoRepo {
	return &MemoRepo{db: sync.Map{}}
}

// Add provides a functionality to save a slice of the ShortURL data into the in-memo repository.
// Since it doesn't depend on any additional readers, it returns the copied value of the slice.
func (m *MemoRepo) Add(batch []ShortURL) ([]ShortURL, error) {
	for _, sURL := range batch {
		m.db.Store(sURL.ID, sURL)
	}

	res := make([]ShortURL, len(batch))
	copy(res, batch)
	return res, nil
}

// Has checks if the repository contains the ShortURL with a specific ID.
func (m *MemoRepo) Has(id string) (bool, error) {
	if _, ok := m.db.Load(id); ok {
		return true, nil
	}

	return false, nil
}

// Get returns the ShortURL value by its ID.
// If the value is missing from the repository, the error will be returned.
func (m *MemoRepo) Get(id string) (ShortURL, error) {
	if sURL, ok := m.db.Load(id); ok {
		return sURL.(ShortURL), nil
	}

	return ShortURL{}, errors.New(apperrors.URLNotFound)
}

// GetAll returns all the ShortURL values created by the specified user.
// If the repository doesn't have any associated value, the empty slice will be returned.
func (m *MemoRepo) GetAll(userID string) ([]ShortURL, error) {
	urls := make([]ShortURL, 0)

	m.db.Range(func(_, v interface{}) bool {
		sURL := v.(ShortURL)
		if sURL.UID == userID {
			urls = append(urls, sURL)
		}
		return true
	})

	return urls, nil
}

// Clear marks all existing values in the repository as deleted.
func (m *MemoRepo) Clear() {
	m.db.Range(func(key, _ interface{}) bool {
		m.db.Delete(key)
		return true
	})
}

// Ping functionality is not supported by the in-memo repository, so this function always return true.
func (m *MemoRepo) Ping() bool {
	return true
}

// Delete marks all specified ShortURL values in repository as deleted.
// The deletion of the value is available only for its owner. All other values will be skipped.
func (m *MemoRepo) Delete(batch []ShortURL) error {
	for _, sURL := range batch {
		stored, ok := m.db.Load(sURL.ID)
		if !ok || stored.(ShortURL).UID != sURL.UID {
			continue
		}

		newURL := stored.(ShortURL)
		newURL.Deleted = true
		m.db.Store(sURL.ID, newURL)
	}

	return nil
}

func (m *MemoRepo) Close() error {
	return nil
}

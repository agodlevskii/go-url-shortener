package storage

import (
	"errors"
	"go-url-shortener/internal/apperrors"
	"sync"
)

type MemoRepo struct {
	db sync.Map
}

func NewMemoryRepo() *MemoRepo {
	return &MemoRepo{db: sync.Map{}}
}

func (m *MemoRepo) Add(batch []ShortURL) ([]ShortURL, error) {
	for _, sURL := range batch {
		m.db.Store(sURL.ID, sURL)
	}

	res := make([]ShortURL, len(batch))
	copy(res, batch)
	return res, nil
}

func (m *MemoRepo) Has(id string) (bool, error) {
	if _, ok := m.db.Load(id); ok {
		return true, nil
	}

	return false, nil
}

func (m *MemoRepo) Get(id string) (ShortURL, error) {
	if sURL, ok := m.db.Load(id); ok {
		return sURL.(ShortURL), nil
	}

	return ShortURL{}, errors.New(apperrors.URLNotFound)
}

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

func (m *MemoRepo) Clear() {
	m.db.Range(func(key, _ interface{}) bool {
		m.db.Delete(key)
		return true
	})
}

func (m *MemoRepo) Ping() bool {
	return true
}

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

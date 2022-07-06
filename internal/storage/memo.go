package storage

import (
	"errors"
	"sync"
)

type MemoRepo struct {
	db sync.Map
}

func NewMemoryRepo() *MemoRepo {
	var db sync.Map
	return &MemoRepo{db: db}
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

func (m *MemoRepo) Get(id string) (string, error) {
	if sURL, ok := m.db.Load(id); ok {
		return sURL.(ShortURL).URL, nil
	}

	return "", errors.New("no matching URL found")
}

func (m *MemoRepo) GetAll(userID string) ([]ShortURL, error) {
	urls := make([]ShortURL, 0)

	m.db.Range(func(_, v any) bool {
		sURL := v.(ShortURL)
		if sURL.UID == userID {
			urls = append(urls, sURL)
		}
		return true
	})

	return urls, nil
}

func (m *MemoRepo) Clear() {
	m.db.Range(func(key, _ any) bool {
		m.db.Delete(key)
		return true
	})
}

func (m *MemoRepo) Ping() bool {
	return true
}

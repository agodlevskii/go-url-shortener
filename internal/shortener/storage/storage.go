package storage

import (
	"errors"
)

type Storager interface {
	Add(id, val string) error
	Has(string) bool
	Get(string) (string, error)
	Remove(string) error
	Clear()
}

type MemoRepo struct {
	db map[string]string
}

func NewMemoryRepo(data map[string]string) MemoRepo {
	repo := MemoRepo{db: make(map[string]string)}
	repo.Init(data)
	return repo
}

func (m MemoRepo) Add(id string, url string) error {
	m.db[id] = url
	return nil
}

func (m MemoRepo) Has(id string) bool {
	return m.db[id] != ""
}

func (m MemoRepo) Get(id string) (string, error) {
	if m.Has(id) {
		return m.db[id], nil
	}

	return "", errors.New("no matching URL found")
}

func (m MemoRepo) Remove(id string) error {
	if m.Has(id) {
		delete(m.db, id)
		return nil
	}

	return errors.New("no matching URL found")
}

func (m MemoRepo) Init(data map[string]string) error {
	for k, v := range data {
		m.db[k] = v
	}

	return nil
}

func (m MemoRepo) Clear() {
	for id := range m.db {
		delete(m.db, id)
	}
}

func AddURLToStorage(repo Storager, id string, url string) error {
	return repo.Add(id, url)
}

func GetURLFromStorage(repo Storager, id string) (string, error) {
	return repo.Get(id)
}

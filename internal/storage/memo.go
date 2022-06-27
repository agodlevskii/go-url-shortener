package storage

import "errors"

type MemoRepo struct {
	db map[string]URLRes
}

func NewMemoryRepo() MemoRepo {
	return MemoRepo{db: make(map[string]URLRes)}
}

func (m MemoRepo) Add(userID, id, url string) error {
	m.db[id] = URLRes{
		url: url,
		uid: userID,
	}

	return nil
}

func (m MemoRepo) Has(id string) (bool, error) {
	if _, ok := m.db[id]; ok {
		return true, nil
	}

	return false, nil
}

func (m MemoRepo) Get(id string) (string, error) {
	if res, ok := m.db[id]; ok {
		return res.url, nil
	}

	return "", errors.New("no matching URL found")
}

func (m MemoRepo) GetAll(userID string) (map[string]string, error) {
	urls := make(map[string]string)
	for k, v := range m.db {
		if v.uid == userID {
			urls[k] = v.url
		}
	}

	return urls, nil
}

func (m MemoRepo) Clear() {
	for id := range m.db {
		delete(m.db, id)
	}
}

package storage

import (
	"bufio"
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type URLRes struct {
	url string
	uid string
}

type Storager interface {
	Add(userID, id, val string) error
	Has(id string) (bool, error)
	Get(id string) (string, error)
	GetAll(userID string) (map[string]string, error)
	Clear()
}

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

type FileRepo struct {
	filename string
}

func NewFileRepo(filename string) (FileRepo, error) {
	if filename == "" {
		return FileRepo{}, errors.New("the filename is missing")
	}

	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return FileRepo{}, err
	}

	if err = file.Close(); err != nil {
		return FileRepo{}, err
	}

	return FileRepo{filename: filename}, nil
}

func (f FileRepo) Add(userID, id, url string) error {
	file, err := os.OpenFile(f.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	if _, err = file.WriteString(id + " : " + url + " : " + userID + "\n"); err != nil {
		return err
	}
	return file.Close()
}

func (f FileRepo) Get(id string) (string, error) {
	file, err := os.OpenFile(f.filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	counter := 0

	for ; scanner.Scan(); counter++ {
		data := strings.Split(scanner.Text(), " : ")

		if len(data) < 3 {
			return "", errors.New("malformed file: " + file.Name())
		}

		if data[0] == id {
			return data[1], nil
		}
	}

	if counter == 0 {
		return "", scanner.Err()
	}

	return "", errors.New("no matching URL found")
}

func (f FileRepo) GetAll(userID string) (map[string]string, error) {
	file, err := os.OpenFile(f.filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Error(err)
		return nil, nil
	}
	defer file.Close()

	urls := make(map[string]string)
	scanner := bufio.NewScanner(file)
	counter := 0

	for ; scanner.Scan(); counter++ {
		data := strings.Split(scanner.Text(), " : ")

		if len(data) < 3 {
			return nil, errors.New("malformed file: " + file.Name())
		}

		if data[2] == userID {
			urls[data[0]] = data[1]
		}
	}

	if counter == 0 {
		return urls, scanner.Err()
	}

	return urls, nil
}

func (f FileRepo) Has(id string) (bool, error) {
	file, err := os.OpenFile(f.filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return false, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return false, err
	}

	if stat.Size() == 0 {
		return false, nil
	}

	scanner := bufio.NewScanner(file)
	counter := 0

	for ; scanner.Scan(); counter++ {
		data := strings.Split(scanner.Text(), " : ")
		if len(data) < 3 {
			return false, nil
		}

		if data[0] == id {
			return true, nil
		}
	}

	return false, nil
}

func (f FileRepo) Clear() {
	if _, err := os.Create(f.filename); err != nil {
		log.Error(err)
	}
}

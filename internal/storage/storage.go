package storage

import (
	"bufio"
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Storager interface {
	Add(id, val string) error
	Has(string) bool
	Get(string) (string, error)
	Clear()
}

type MemoRepo struct {
	db map[string]string
}

func NewMemoryRepo() MemoRepo {
	return MemoRepo{db: make(map[string]string)}
}

func (m MemoRepo) Add(id string, url string) error {
	m.db[id] = url
	return nil
}

func (m MemoRepo) Has(id string) bool {
	_, ok := m.db[id]
	return ok
}

func (m MemoRepo) Get(id string) (string, error) {
	if m.Has(id) {
		return m.db[id], nil
	}

	return "", errors.New("no matching URL found")
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

func (f FileRepo) Add(id string, url string) error {
	file, err := os.OpenFile(f.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	if _, err = file.WriteString(id + " : " + url + "\n"); err != nil {
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

		if len(data) < 2 {
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

func (f FileRepo) Has(id string) bool {
	file, err := os.OpenFile(f.filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Error(err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	counter := 0

	for ; scanner.Scan(); counter++ {
		data := strings.Split(scanner.Text(), " : ")
		if len(data) < 2 {
			return false
		}

		if data[0] == id {
			return true
		}
	}

	return false
}

func (f FileRepo) Clear() {
	if _, err := os.Create(f.filename); err != nil {
		log.Error(err)
	}
}

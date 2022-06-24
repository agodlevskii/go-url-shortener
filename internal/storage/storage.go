package storage

import (
	"bufio"
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Storager interface {
	Add(userId, id, val string) error
	Has(userId, id string) (bool, error)
	Get(userId, id string) (string, error)
	GetAll(userId string) (map[string]string, error)
	Clear()
}

type MemoRepo struct {
	db map[string]map[string]string
}

func NewMemoryRepo() MemoRepo {
	return MemoRepo{db: make(map[string]map[string]string)}
}

func (m MemoRepo) Add(userId, id, url string) error {
	us, ok := m.db[userId]
	if !ok {
		m.db[userId] = make(map[string]string)
		us = m.db[userId]
	}

	us[id] = url
	return nil
}

func (m MemoRepo) Has(userId, id string) (bool, error) {
	if us, ok := m.db[userId]; ok {
		if _, ok := us[id]; ok {
			return true, nil
		}
	}

	return false, nil
}

func (m MemoRepo) Get(userId, id string) (string, error) {
	if us, ok := m.db[userId]; ok {
		if url, ok := us[id]; ok {
			return url, nil
		}
	}

	return "", errors.New("no matching URL found")
}

func (m MemoRepo) GetAll(userId string) (map[string]string, error) {
	if us, ok := m.db[userId]; ok {
		return us, nil
	}

	return nil, nil
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

func (f FileRepo) Add(userId, id, url string) error {
	file, err := os.OpenFile(f.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	if _, err = file.WriteString(userId + " : " + id + " : " + url + "\n"); err != nil {
		return err
	}
	return file.Close()
}

func (f FileRepo) Get(userId, id string) (string, error) {
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

		if data[0] == userId && data[1] == id {
			return data[2], nil
		}
	}

	if counter == 0 {
		return "", scanner.Err()
	}

	return "", errors.New("no matching URL found")
}

func (f FileRepo) GetAll(userId string) (map[string]string, error) {
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

		if data[0] == userId {
			urls[data[1]] = data[2]
		}
	}

	if counter == 0 {
		return urls, scanner.Err()
	}

	return urls, nil
}

func (f FileRepo) Has(userId, id string) (bool, error) {
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

		if data[0] == userId && data[1] == id {
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

package storage

import (
	"bufio"
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

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

func (f FileRepo) Add(userID string, batch map[string]string) (map[string]string, error) {
	file, err := os.OpenFile(f.filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	res := make(map[string]string, len(batch))
	for id, url := range batch {
		if _, err = file.WriteString(id + " : " + url + " : " + userID + "\n"); err != nil {
			return nil, err
		}

		res[url] = id
	}

	return res, file.Close()
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
	if err := os.Remove(f.filename); err != nil {
		log.Error(err)
	}
}

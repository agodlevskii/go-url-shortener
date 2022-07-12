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

func (f FileRepo) Add(batch []ShortURL) ([]ShortURL, error) {
	file, err := os.OpenFile(f.filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	w := bufio.NewWriter(file)
	for _, sURL := range batch {
		_, err = w.WriteString(sURL.ID + " : " + sURL.URL + " : " + sURL.UID + "false\n")
		if err != nil {
			return nil, err
		}
	}

	if err = w.Flush(); err != nil {
		return nil, err
	}

	res := make([]ShortURL, len(batch))
	copy(res, batch)
	return res, file.Close()
}

func (f FileRepo) Get(id string) (ShortURL, error) {
	file, err := os.OpenFile(f.filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return ShortURL{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		data := strings.Split(text, " : ")

		if !isEntryValid(data) {
			return ShortURL{}, errors.New("malformed file: " + file.Name())
		}

		if data[0] == id {
			return ShortURL{
				ID:      data[0],
				URL:     data[1],
				UID:     data[2],
				Deleted: data[3] == "true",
			}, nil
		}
	}

	if scanner.Err() != nil {
		return ShortURL{}, scanner.Err()
	}

	return ShortURL{}, errors.New("no matching URL found: " + id)
}

func (f FileRepo) GetAll(userID string) ([]ShortURL, error) {
	file, err := os.OpenFile(f.filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Error(err)
		return nil, nil
	}
	defer file.Close()

	urls := make([]ShortURL, 0)
	scanner := bufio.NewScanner(file)
	counter := 0

	for ; scanner.Scan(); counter++ {
		data := strings.Split(scanner.Text(), " : ")

		if !isEntryValid(data) {
			return nil, errors.New("malformed file: " + file.Name())
		}

		if data[2] == userID {
			urls = append(urls, ShortURL{
				ID:      data[0],
				URL:     data[1],
				UID:     userID,
				Deleted: data[3] == "true",
			})
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
		if !isEntryValid(data) {
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

func (f FileRepo) Ping() bool {
	_, err := os.Stat(f.filename)
	return err == nil
}

func (f FileRepo) Delete(batch []ShortURL) error {
	file, err := os.OpenFile(f.filename, os.O_RDONLY, 0777)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer file.Close()

	IDtoSURL := make(map[string]ShortURL, len(batch))
	for _, v := range batch {
		IDtoSURL[v.ID] = v
	}

	restore := make([]ShortURL, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), " : ")

		if !isEntryValid(data) {
			return errors.New("malformed file: " + file.Name())
		}

		if sURL := IDtoSURL[data[0]]; sURL.UID == "" {
			restore = append(restore, ShortURL{
				ID:      data[0],
				URL:     data[1],
				UID:     data[2],
				Deleted: sURL.UID == data[2],
			})
		}
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	f.Clear()
	if _, err = f.Add(restore); err != nil {
		return err
	}

	return nil
}

func isEntryValid(entry []string) bool {
	return len(entry) == 4
}

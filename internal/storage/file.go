package storage

import (
	"bufio"
	"errors"
	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"
	"go-url-shortener/internal/apperrors"
	"os"
)

type FileRepo struct {
	filename string
}

func NewFileRepo(filename string) (FileRepo, error) {
	if filename == "" {
		return FileRepo{}, errors.New(apperrors.FilenameMissing)
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return FileRepo{}, err
	}

	return FileRepo{filename: filename}, file.Close()
}

func (f FileRepo) Add(batch []ShortURL) ([]ShortURL, error) {
	file, err := os.OpenFile(f.filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, sURL := range batch {
		if _, err = w.WriteString(ShortURLToRepoString(sURL)); err != nil {
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
		if sURL, err := RepoStringToShortURL(scanner.Text()); err != nil || sURL.ID == id {
			return sURL, err
		}
	}

	if scanner.Err() != nil {
		return ShortURL{}, scanner.Err()
	}

	return ShortURL{}, errors.New(pretty.Sprintf("%s: %s", apperrors.URLNotFound, id))
}

func (f FileRepo) GetAll(userID string) ([]ShortURL, error) {
	file, err := os.OpenFile(f.filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	urls := make([]ShortURL, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sURL, err := RepoStringToShortURL(scanner.Text())
		if err != nil {
			return nil, err
		}

		if sURL.UID == userID {
			urls = append(urls, sURL)
		}
	}

	return urls, scanner.Err()
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
	for scanner.Scan() {
		sURL, err := RepoStringToShortURL(scanner.Text())
		if err != nil {
			return false, err
		}

		if sURL.ID == id {
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
		stored, err := RepoStringToShortURL(scanner.Text())
		if err != nil {
			return err
		}

		stored.Deleted = IDtoSURL[stored.ID].UID != ""
		restore = append(restore, stored)
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

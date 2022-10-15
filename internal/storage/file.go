package storage

import (
	"bufio"
	"errors"
	"os"

	"go-url-shortener/internal/apperrors"

	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"
)

// FileRepo describes the file-based implementation of the Storager interface.
type FileRepo struct {
	filename string
}

// NewFileRepo returns a new instance of the FileRepo type.
// If the filename is missing, the error will be returned.
// If the file with the associated filename is missing, it will be created.
// Otherwise, its content will be removed.
func NewFileRepo(filename string) (FileRepo, error) {
	if filename == "" {
		return FileRepo{}, errors.New(apperrors.FilenameMissing)
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC, 0o777)
	if err != nil {
		return FileRepo{}, err
	}

	return FileRepo{filename: filename}, file.Close()
}

// Add provides a functionality to save a slice of the ShortURL data into the file-based repository.
// If the file with the associated filename is missing, it will be created.
// Otherwise, it will be opened for writing.
func (f FileRepo) Add(batch []ShortURL) ([]ShortURL, error) {
	file, err := os.OpenFile(f.filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o777)
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

// Get returns the ShortURL value by its ID.
// If the value is missing from the repository, the error will be returned.
// If the file with the associated filename is missing, it will be created.
// Otherwise, it will be opened for reading.
func (f FileRepo) Get(id string) (ShortURL, error) {
	file, err := os.OpenFile(f.filename, os.O_RDONLY|os.O_CREATE, 0o777)
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

// GetAll returns all the ShortURL values created by the specified user.
// If the repository doesn't have any associated value, the empty slice will be returned.
// Otherwise, it will be opened for reading.
func (f FileRepo) GetAll(userID string) ([]ShortURL, error) {
	file, err := os.OpenFile(f.filename, os.O_RDONLY|os.O_CREATE, 0o777)
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

// Has checks if the repository contains the ShortURL with a specific ID.
// If the file with the associated filename is missing, it will be created.
// Otherwise, it will be opened for reading.
func (f FileRepo) Has(id string) (bool, error) {
	file, err := os.OpenFile(f.filename, os.O_RDONLY|os.O_CREATE, 0o777)
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

// Clear removes the associated file from the hard drive.
func (f FileRepo) Clear() {
	if err := os.Remove(f.filename); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Error(err)
	}
}

// Ping checks if the associated file exists.
func (f FileRepo) Ping() bool {
	_, err := os.Stat(f.filename)
	return err == nil
}

// Delete marks all specified ShortURL values in repository as deleted.
// The deletion of the value is available only for its owner. All other values will be skipped.
// If the file with the associated filename is missing, it will be created.
// Otherwise, it will be opened for reading.
func (f FileRepo) Delete(batch []ShortURL) error {
	file, err := os.OpenFile(f.filename, os.O_RDONLY, 0o777)
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

func (f FileRepo) Close() error {
	return nil
}

// Package storage includes the interfaces and functions related to the storing functionality.
package storage

import (
	"errors"
	"go-url-shortener/internal/apperrors"
	"strconv"
	"strings"
)

// ShortURL describes the type of data stored in the entities that implement the Storager interface.
type ShortURL struct {
	ID      string
	URL     string
	UID     string
	Deleted bool
}

// Storager describes the functionality that can be performed on the storage instance.
type Storager interface {
	Add(batch []ShortURL) ([]ShortURL, error)
	Clear()
	Delete(batch []ShortURL) error
	Get(id string) (ShortURL, error)
	GetAll(userID string) ([]ShortURL, error)
	Has(id string) (bool, error)
	Ping() bool
	Close() error
}

// RepoStrSep describes the string that separates the ShortURL field values in the file-based Storager implementation.
const RepoStrSep = " : "

// ShortURLToRepoString converts the ShortURL instance into a string for the file-based Storager interface.
// It uses the RepoStrSep constant to divide the field values.
func ShortURLToRepoString(sURL ShortURL) string {
	return sURL.ID + RepoStrSep + sURL.URL + RepoStrSep + sURL.UID + RepoStrSep + strconv.FormatBool(sURL.Deleted) + "\n"
}

// RepoStringToShortURL converts a string for the file-based Storager interface into the ShortURL instance.
// It uses the RepoStrSep constant to divide the field values.
func RepoStringToShortURL(str string) (ShortURL, error) {
	entry := strings.Split(str, RepoStrSep)
	if !isEntryValid(entry) {
		return ShortURL{}, errors.New(apperrors.RepoEntryInvalid)
	}

	return ShortURL{
		ID:      entry[0],
		URL:     entry[1],
		UID:     entry[2],
		Deleted: entry[3] == "true",
	}, nil
}

// isEntryValid validates the string for the file-based Storager, so it would include all ShortURL fields.
func isEntryValid(entry []string) bool {
	return len(entry) == 4
}

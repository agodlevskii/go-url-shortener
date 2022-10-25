// Package storage includes the interfaces and functions related to the storing functionality.
package storage

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"go-url-shortener/internal/apperrors"
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
	Add(ctx context.Context, batch []ShortURL) ([]ShortURL, error)
	Clear(ctx context.Context)
	Delete(ctx context.Context, batch []ShortURL) error
	Get(ctx context.Context, id string) (ShortURL, error)
	GetAll(ctx context.Context, userID string) ([]ShortURL, error)
	Has(ctx context.Context, id string) (bool, error)
	Ping(ctx context.Context) bool
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

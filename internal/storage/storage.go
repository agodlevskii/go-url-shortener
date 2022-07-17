package storage

import (
	"errors"
	"go-url-shortener/internal/apperrors"
	"strconv"
	"strings"
)

type ShortURL struct {
	ID      string
	URL     string
	UID     string
	Deleted bool
}

type Storager interface {
	Add(batch []ShortURL) ([]ShortURL, error)
	Clear()
	Delete(batch []ShortURL) error
	Get(id string) (ShortURL, error)
	GetAll(userID string) ([]ShortURL, error)
	Has(id string) (bool, error)
	Ping() bool
}

const RepoStrSep = " : "

func ShortURLToRepoString(sURL ShortURL) string {
	return sURL.ID + RepoStrSep + sURL.URL + RepoStrSep + sURL.UID + RepoStrSep + strconv.FormatBool(sURL.Deleted) + "\n"
}

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

func isEntryValid(entry []string) bool {
	return len(entry) == 4
}

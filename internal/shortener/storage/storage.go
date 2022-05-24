package storage

import (
	"errors"
	"strings"
)

var Storage = make(map[string]string)

func AddURLToStorage(url string) string {
	if url == "" {
		return ""
	}

	parts := strings.Split(url, "://")
	urlToShorten := url
	if len(parts) > 1 {
		urlToShorten = parts[1]
	}

	surl := urlToShorten[:len(urlToShorten)/2]
	Storage[surl] = url
	return surl
}

func GetURLFromStorage(id string) (string, error) {
	url := Storage[id]
	if url == "" {
		return "", errors.New("the URL with associated ID is not found")
	}

	return url, nil
}

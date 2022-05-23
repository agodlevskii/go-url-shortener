package shortener

import (
	"errors"
	"strings"
)

var storage = make(map[string]string)

func AddUrlToStorage(url string) string {
	parts := strings.Split(url, "://")
	urlToShorten := url
	if len(parts) > 1 {
		urlToShorten = parts[1]
	}

	surl := urlToShorten[:len(urlToShorten)/2]
	storage[surl] = url
	return surl
}

func GetUrlFromStorage(id string) (string, error) {
	url := storage[id]
	if url == "" {
		return "", errors.New("the URL with associated ID is not found")
	}

	return url, nil
}

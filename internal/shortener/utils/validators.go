package utils

import "net/url"

func IsURLValid(rawURL *url.URL) bool {
	return IsURLStringValid(rawURL.Path)
}

func IsURLStringValid(rawURL string) bool {
	_, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	return true
}

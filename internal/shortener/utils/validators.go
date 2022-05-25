package utils

import "net/url"

func IsURLValid(rawURL *url.URL) bool {
	return IsURLStringValid(rawURL.Path)
}

func IsURLStringValid(rawURL string) bool {
	_, err := url.ParseRequestURI(rawURL)
	return err == nil
}

package validators

import "net/url"

func IsURLStringValid(rawURL string) bool {
	_, err := url.ParseRequestURI(rawURL)
	return err == nil
}

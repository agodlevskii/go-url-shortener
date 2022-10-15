// Package validators provides the functionality to validate the application data.
package validators

import (
	"net/url"

	log "github.com/sirupsen/logrus"
)

// IsURLStringValid checks if the URL string is of a valid format.
// If any issues arise during the URL parsing, the false value will be returned.
func IsURLStringValid(rawURL string) bool {
	if _, err := url.ParseRequestURI(rawURL); err != nil {
		log.Error(err)
		return false
	}

	return true
}

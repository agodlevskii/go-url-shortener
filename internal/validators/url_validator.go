package validators

import (
	"net/url"

	log "github.com/sirupsen/logrus"
)

func IsURLStringValid(rawURL string) bool {
	if _, err := url.ParseRequestURI(rawURL); err != nil {
		log.Error(err)
		return false
	}

	return true
}

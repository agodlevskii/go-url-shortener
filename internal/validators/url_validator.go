package validators

import (
	log "github.com/sirupsen/logrus"
	"net/url"
)

func IsURLStringValid(rawURL string) bool {
	if _, err := url.ParseRequestURI(rawURL); err != nil {
		log.Error(err)
		return false
	}

	return true
}

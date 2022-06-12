package validators

import (
	log "github.com/sirupsen/logrus"
	"net/url"
)

func IsURLStringValid(rawURL string) bool {
	_, err := url.ParseRequestURI(rawURL)
	if err != nil {
		log.Error(err)
	}

	return err == nil
}

package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// GetHomePage handles the request for the index page.
// The map of the templates is being passed from the main handler.
// If the required template is missing from the map or malformed, the user gets an error response.
func GetHomePage(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("The URL shortener is up and running.")); err != nil {
		log.Error(err)
	}
}

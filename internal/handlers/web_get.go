package handlers

import (
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"go-url-shortener/internal/storage"
	"net/http"
)

func GetFullURL(db storage.Storager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		sURL, err := db.Get(id)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if sURL.Deleted {
			http.Error(w, "The requested URL is not available.", http.StatusGone)
		}

		http.Redirect(w, r, sURL.URL, http.StatusTemporaryRedirect)
	}
}

package handlers

import (
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"go-url-shortener/internal/middlewares"
	"go-url-shortener/internal/storage"
	"net/http"
)

func GetFullURL(db storage.Storager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := middlewares.GetUserId(r)
		if err != nil {
			log.Error(err)
			http.Error(w, "Cannot identify a user", http.StatusInternalServerError)
			return
		}

		id := chi.URLParam(r, "id")
		url, err := db.Get(userId, id)
		log.Info("userId: ", userId)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go-url-shortener/internal/middlewares"
	"go-url-shortener/internal/storage"
	"net/http"
)

func DeleteShortURLs(db storage.Storager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middlewares.GetUserID(r)
		if err != nil {
			log.Error(err)
			http.Error(w, "couldn't identify the user", http.StatusInternalServerError)
			return
		}

		var req []string
		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error(err)
			http.Error(w, "Couldn't delete the records", http.StatusInternalServerError)
			return
		}

		batch := make([]storage.ShortURL, len(req))
		for i, v := range req {
			batch[i] = storage.ShortURL{
				ID:  v,
				UID: userID,
			}
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(""))
	}
}

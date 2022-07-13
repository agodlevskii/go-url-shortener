package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go-url-shortener/internal/middlewares"
	"go-url-shortener/internal/storage"
	"net/http"
)

func DeleteShortURLs(db storage.Storager, poolSize int) func(w http.ResponseWriter, r *http.Request) {
	pool := make(chan func(), poolSize)
	for i := 0; i < poolSize; i++ {
		go func() {
			for f := range pool {
				f()
			}
		}()
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middlewares.GetUserID(r)
		if err != nil {
			log.Error(err)
			http.Error(w, "couldn't identify the user", http.StatusBadRequest)
			return
		}

		var ids []string
		if err = json.NewDecoder(r.Body).Decode(&ids); err != nil || len(ids) == 0 {
			log.Error(err)
			http.Error(w, "Please provide a valid list of IDs.", http.StatusBadRequest)
			return
		}

		go func() {
			pool <- func() {
				deleteURLs(db, userID, ids)
			}
		}()

		w.WriteHeader(http.StatusAccepted)
	}
}

func deleteURLs(db storage.Storager, userID string, ids []string) {
	batch := make([]storage.ShortURL, len(ids))
	for i, v := range ids {
		batch[i] = storage.ShortURL{
			ID:  v,
			UID: userID,
		}
	}

	err := db.Delete(batch)
	if err != nil {
		log.Error(err)
	}
}

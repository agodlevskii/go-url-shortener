package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/middlewares"
	"go-url-shortener/internal/storage"
	"net/http"
)

func UserURLsHandler(db storage.Storager, baseURL string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middlewares.GetUserID(r)
		if err != nil {
			apperrors.HandleUserError(w)
			return
		}

		list := getUserLinks(db, userID, baseURL)
		if len(list) == 0 {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("No results found."))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(list); err != nil {
			apperrors.HandleInternalError(w)
		}
	}
}

func getUserLinks(db storage.Storager, userID, baseURL string) []UserLink {
	urls, err := db.GetAll(userID)
	if err != nil {
		log.Error(err)
		return nil
	}

	if len(urls) == 0 {
		return nil
	}

	links := make([]UserLink, 0)
	for _, url := range urls {
		links = append(links, UserLink{
			Short:    baseURL + "/" + url.ID,
			Original: url.URL,
		})
	}

	return links
}

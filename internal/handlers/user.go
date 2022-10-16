package handlers

import (
	"encoding/json"
	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/middlewares"
	"go-url-shortener/internal/storage"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// GetUserLinks returns the list of the user-associated links.
// The user is being identified based on a request cookie.
// The response includes full information on the stored link, including the deletion flag.
func GetUserLinks(db storage.Storager, baseURL string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middlewares.GetUserID(r)
		if err != nil {
			apperrors.HandleUserError(w)
			return
		}

		list := getLinks(db, userID, baseURL)
		if len(list) == 0 {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNoContent)
			if _, err = w.Write([]byte("No results found.")); err != nil {
				log.Error(err)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(list); err != nil {
			apperrors.HandleInternalError(w)
		}
	}
}

// DeleteUserLinks deletes the specified entities from the list of the user-associated links.
// The user is being identified based on a request cookie.
// The links must be passed as an array of strings in the request body.
// The handler doesn't remove the links, but validates the request and marks the passed entities for deletion.
func DeleteUserLinks(db storage.Storager, poolSize int) func(w http.ResponseWriter, r *http.Request) {
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
			apperrors.HandleUserError(w)
			return
		}

		var ids []string
		if err = json.NewDecoder(r.Body).Decode(&ids); err != nil || len(ids) == 0 {
			apperrors.HandleHTTPError(w, apperrors.NewError(apperrors.IDsListFormat, err), http.StatusBadRequest)
			return
		}

		go func() {
			pool <- func() {
				deleteLinks(db, userID, ids)
			}
		}()

		w.WriteHeader(http.StatusAccepted)
	}
}

// getLinks recovers the user-associated links from the repository.
// In case if there are no links to return, the function returns nil instead of the empty slice.
func getLinks(db storage.Storager, userID, baseURL string) []UserLink {
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

// getLinks deletes the user-associated links from the repository.
// The listed entities remain in the repository, but each of them gets their deletion flag set to true.
func deleteLinks(db storage.Storager, userID string, ids []string) {
	batch := make([]storage.ShortURL, len(ids))
	for i, v := range ids {
		batch[i] = storage.ShortURL{
			ID:  v,
			UID: userID,
		}
	}

	if err := db.Delete(batch); err != nil {
		log.Error(err)
	}
}

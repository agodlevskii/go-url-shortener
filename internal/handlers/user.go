package handlers

import (
	"context"
	"encoding/json"
	"go-url-shortener/internal/services"
	"net/http"

	log "github.com/sirupsen/logrus"

	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/middlewares"
	"go-url-shortener/internal/storage"
)

// GetUserLinks returns the list of the user-associated links.
// The user is being identified based on a request cookie.
// The response includes full information on the stored link, including the deletion flag.
func GetUserLinks(db storage.Storager, cfg APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middlewares.GetUserID(cfg, r)
		if err != nil {
			apperrors.HandleUserError(w)
			return
		}

		list, err := services.GetUserURLs(r.Context(), db, userID, cfg.GetBaseURL())
		if err != nil {
			apperrors.HandleHTTPError(w, apperrors.NewError("", err), http.StatusInternalServerError)
			return
		}

		if len(list) == 0 {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNoContent)
			if _, err = w.Write([]byte("No results found.")); err != nil {
				apperrors.HandleHTTPError(w, apperrors.NewError("", err), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(list); err != nil {
			apperrors.HandleHTTPError(w, apperrors.NewError("", err), http.StatusInternalServerError)
		}
	}
}

// DeleteUserLinks deletes the specified entities from the list of the user-associated links.
// The user is being identified based on a request cookie.
// The links must be passed as an array of strings in the request body.
// The handler doesn't remove the links, but validates the request and marks the passed entities for deletion.
func DeleteUserLinks(db storage.Storager, cfg APIConfig) http.HandlerFunc {
	ps := cfg.GetPoolSize()
	pool := make(chan func(), ps)
	for i := 0; i < ps; i++ {
		go func() {
			for f := range pool {
				f()
			}
		}()
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middlewares.GetUserID(cfg, r)
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
				if err = services.DeleteUserURLs(context.Background(), db, userID, ids); err != nil {
					log.Error(err)
				}
			}
		}()

		w.WriteHeader(http.StatusAccepted)
	}
}

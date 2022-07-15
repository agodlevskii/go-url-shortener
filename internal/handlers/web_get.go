package handlers

import (
	"github.com/go-chi/chi/v5"
	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/storage"
	"net/http"
)

func GetFullURL(db storage.Storager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		sURL, err := db.Get(id)
		if err != nil {
			apperrors.HandleHTTPError(w, apperrors.NewError("", err), http.StatusBadRequest)
			return
		}

		if sURL.Deleted {
			apperrors.HandleHTTPError(w, apperrors.NewError(apperrors.URLGone, nil), http.StatusGone)
			return
		}

		http.Redirect(w, r, sURL.URL, http.StatusTemporaryRedirect)
	}
}

package handlers

import (
	"go-url-shortener/internal/services"
	"net/http"

	_ "github.com/jackc/pgx/v4" // SQL driver
	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/storage"
)

// Ping handles the DB status request.
// If the ping fails, the user gets the error response.
func Ping(db storage.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ping := services.Ping(r.Context(), db); !ping {
			apperrors.HandleInternalError(w)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("DB is up and running")); err != nil {
			apperrors.HandleHTTPError(w, apperrors.NewError("", err), http.StatusInternalServerError)
		}
	}
}

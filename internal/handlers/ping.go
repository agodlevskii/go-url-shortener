package handlers

import (
	"net/http"

	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/storage"

	log "github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v4" // SQL driver
)

// Ping handles the DB status request.
// If the ping fails, the user gets the error response.
func Ping(db storage.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ping := db.Ping(r.Context()); !ping {
			apperrors.HandleInternalError(w)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("DB is up and running")); err != nil {
			log.Error(err)
		}
	}
}

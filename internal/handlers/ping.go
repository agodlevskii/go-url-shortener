package handlers

import (
	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/storage"
	"net/http"

	_ "github.com/jackc/pgx/v4"
)

func Ping(db storage.Storager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if ping := db.Ping(); !ping {
			apperrors.HandleInternalError(w)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("DB is up and running"))
	}
}

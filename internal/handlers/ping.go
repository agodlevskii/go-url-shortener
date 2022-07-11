package handlers

import (
	_ "github.com/jackc/pgx/v4"
	"go-url-shortener/internal/storage"
	"net/http"
)

func Ping(db storage.Storager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if ping := db.Ping(); !ping {
			http.Error(w, "couldn't connect to DB", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("DB is up and running"))
	}
}

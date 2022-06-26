package handlers

import (
	"context"
	"github.com/jackc/pgx/v4"
	"go-url-shortener/internal"
	"net/http"
)

func Ping() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := pgx.Connect(context.Background(), internal.Config.DBURL)
		if err != nil {
			http.Error(w, "couldn't connect to DB", http.StatusInternalServerError)
			return
		}
		defer conn.Close(context.Background())

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("DB is up and running"))
	}
}

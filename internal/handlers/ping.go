package handlers

import (
	"context"
	"github.com/jackc/pgx/v4"
	"go-url-shortener/internal"
	"net/http"
	"time"
)

func Ping() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		conn, err := pgx.Connect(ctx, internal.Config.DBURL)
		if err != nil {
			http.Error(w, "couldn't connect to DB", http.StatusInternalServerError)
			return
		}
		defer conn.Close(ctx)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("DB is up and running"))
	}
}

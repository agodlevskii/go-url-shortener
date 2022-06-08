package main

import (
	log "github.com/sirupsen/logrus"
	"go-url-shortener/configs"
	"go-url-shortener/internal/handlers"
	"go-url-shortener/internal/storage"
	"net/http"
)

func main() {
	db := storage.NewMemoryRepo()
	r := handlers.NewShortenerRouter(db)

	err := http.ListenAndServe(configs.Host+":"+configs.Port, r)
	if err != nil {
		log.Error(err)
	}
}

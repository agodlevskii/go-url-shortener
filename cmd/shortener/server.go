package main

import (
	log "github.com/sirupsen/logrus"
	"go-url-shortener/configs"
	"go-url-shortener/internal/handlers"
	"net/http"
)

func main() {
	r := handlers.NewShortenerRouter()

	err := http.ListenAndServe(configs.Host+":"+configs.Port, r)
	if err != nil {
		log.Error(err)
	}
}

package main

import (
	log "github.com/sirupsen/logrus"
	"go-url-shortener/internal"
	"go-url-shortener/internal/handlers"
	"go-url-shortener/internal/storage"
	"net/http"
)

func main() {
	err := internal.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := getRepo()
	if err != nil {
		log.Fatal(err)
	}

	r := handlers.NewShortenerRouter(repo, internal.Config.BaseURL)

	err = http.ListenAndServe(internal.Config.Addr, r)
	if err != nil {
		log.Fatal(err)
	}
}

func getRepo() (storage.Storager, error) {
	if internal.Config.Filename == "" {
		return storage.NewMemoryRepo(), nil
	}
	return storage.NewFileRepo(internal.Config.Filename)
}

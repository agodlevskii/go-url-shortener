package main

import (
	log "github.com/sirupsen/logrus"
	"go-url-shortener/internal/config"
	"go-url-shortener/internal/handlers"
	"go-url-shortener/internal/storage"
	"net/http"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := getRepo(cfg)
	if err != nil {
		log.Fatal(err)
	}

	r := handlers.NewShortenerRouter(repo, cfg.BaseURL, 10)
	log.Fatal(http.ListenAndServe(cfg.Addr, r))
}

func getRepo(cfg *config.Config) (storage.Storager, error) {
	if cfg.DBURL != "" {
		return storage.NewDBRepo(cfg.DBURL)
	}
	if cfg.Filename != "" {
		return storage.NewFileRepo(cfg.Filename)
	}
	return storage.NewMemoryRepo(), nil
}

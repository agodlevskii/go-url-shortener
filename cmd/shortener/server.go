package main

import (
	log "github.com/sirupsen/logrus"
	"go-url-shortener/internal/config"
	"go-url-shortener/internal/handlers"
	"go-url-shortener/internal/storage"
	"net/http"
)

func main() {
	cfg := config.GetConfig()
	repo, err := getRepo(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer func(repo storage.Storager) {
		if err := repo.Close(); err != nil {
			log.Error(err)
		}
	}(repo)

	r := handlers.NewShortenerRouter(repo)
	if err = http.ListenAndServe(cfg.Addr, r); err != nil {
		log.Fatal(err)
	}
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

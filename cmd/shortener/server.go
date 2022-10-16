package main

import (
	"go-url-shortener/internal/config"
	"go-url-shortener/internal/handlers"
	"go-url-shortener/internal/storage"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.GetConfig()
	repo, err := getRepo(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer func(repo storage.Storager) {
		if cErr := repo.Close(); cErr != nil {
			log.Error(cErr)
		}
	}(repo)

	r := handlers.NewShortenerRouter(repo)
	if err = getServer(cfg.Addr, r).ListenAndServe(); err != nil {
		log.Error(err)
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

func getServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 3 * time.Second,
	}
}

package main

import (
	"context"
	"net/http"
	"time"

	"go-url-shortener/internal/config"
	"go-url-shortener/internal/handlers"
	"go-url-shortener/internal/storage"

	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.New(config.WithEnv(), config.WithFlags())
	repo, err := getRepo(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer func(repo storage.Storager) {
		if cErr := repo.Close(); cErr != nil {
			log.Error(cErr)
		}
	}(repo)

	r := handlers.NewShortenerRouter(cfg, repo)
	if err = getServer(cfg.Addr, r).ListenAndServe(); err != nil {
		log.Error(err)
	}
}

func getRepo(ctx context.Context, cfg *config.Config) (storage.Storager, error) {
	if cfg.DBURL != "" {
		return storage.NewDBRepo(ctx, cfg.DBURL)
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

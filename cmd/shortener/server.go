package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"go-url-shortener/internal/config"
	"go-url-shortener/internal/handlers"
	"go-url-shortener/internal/storage"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	printCompilationInfo()

	cfg := config.New(config.WithEnv(), config.WithFlags())
	flag.Parse()

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
	serv := getServer(cfg, r)
	defer func(serv *http.Server) {
		if sErr := serv.Close(); sErr != nil {
			log.Error(sErr)
		}
	}(serv)

	if cfg.IsSecure() {
		if err = serv.ListenAndServeTLS("tls.crt", "tls.key"); err != nil {
			log.Error(err)
		}
	} else {
		if err = serv.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}
}

func getRepo(ctx context.Context, cfg *config.Config) (storage.Storager, error) {
	if cfg.GetDBURL() != "" {
		return storage.NewDBRepo(ctx, cfg.GetDBURL())
	}
	if cfg.GetStorageFileName() != "" {
		return storage.NewFileRepo(cfg.GetStorageFileName())
	}
	return storage.NewMemoryRepo(), nil
}

func getServer(cfg *config.Config, handler http.Handler) *http.Server {
	s := &http.Server{
		Addr:              cfg.GetServerAddr(),
		Handler:           handler,
		ReadHeaderTimeout: 3 * time.Second,
	}

	if cfg.IsSecure() {
		log.Info("secure")
		s.TLSConfig = getTLSConfig()
	}

	return s
}

func printCompilationInfo() {
	version := getCompilationInfoValue(buildVersion)
	date := getCompilationInfoValue(buildDate)
	commit := getCompilationInfoValue(buildCommit)
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n\n", version, date, commit)
}

func getCompilationInfoValue(v string) string {
	if v != "" {
		return v
	}
	return "N/A"
}

func getTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
}

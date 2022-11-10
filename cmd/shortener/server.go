package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	cfg := config.New(config.WithEnv(), config.WithFlags(), config.WithFile())

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
	idleConnectionsClosed := make(chan struct{})

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

		<-exit
		stopServer(serv)
		close(idleConnectionsClosed)
	}()

	go startServer(serv, cfg)
	<-idleConnectionsClosed
}

func startServer(s *http.Server, cfg *config.Config) {
	var err error

	if cfg.IsSecure() {
		err = s.ListenAndServeTLS("tls.crt", "tls.key")
	} else {
		err = s.ListenAndServe()
	}

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error(err)
	}
}

func stopServer(s *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error(err)
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
		},
	}
}

package main

import (
	log "github.com/sirupsen/logrus"
	"go-url-shortener/internal/handlers"
	"go-url-shortener/internal/storage"
	"net/http"
	"os"
)

const (
	addrKey         = "SERVER_ADDRESS"
	storageFileName = "FILE_STORAGE_PATH"
)

func main() {
	repo, err := getRepo()
	if err != nil {
		log.Error(err)
	}

	r := handlers.NewShortenerRouter(repo)
	addr, err := getServerAddress()
	if err != nil {
		log.Error(err)
	}

	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Error(err)
	}
}

func getRepo() (storage.Storager, error) {
	os.Setenv(storageFileName, "storage")
	filename := os.Getenv(storageFileName)
	if filename == "" {
		return storage.NewMemoryRepo(), nil
	}

	return storage.NewFileRepo(filename)
}

func getServerAddress() (string, error) {
	var err error
	addr, ok := os.LookupEnv(addrKey)

	if !ok {
		addr = "localhost:8080"
		err = os.Setenv(addrKey, addr)
	}

	return addr, err
}

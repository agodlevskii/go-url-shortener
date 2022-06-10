package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"go-url-shortener/internal/handlers"
	"go-url-shortener/internal/storage"
	"net/http"
	"os"
)

const (
	addrKey         = "SERVER_ADDRESS"
	baseKey         = "BASE_URL"
	storageFileName = "FILE_STORAGE_PATH"
)

var config struct {
	addr     string
	baseURL  string
	filename string
}

func main() {
	flag.Parse()
	repo, err := getRepo()
	if err != nil {
		log.Error(err)
	}

	r := handlers.NewShortenerRouter(repo, config.baseURL)

	err = http.ListenAndServe(config.addr, r)
	if err != nil {
		log.Error(err)
	}
}

func init() {
	flag.StringVar(&config.addr, "a", "", "The application server address")
	flag.StringVar(&config.baseURL, "b", "", "The application server port")
	flag.StringVar(&config.filename, "f", "", "The file storage name")

	setServerAddress()
	setBaseURL()
}

func getRepo() (storage.Storager, error) {
	if config.filename == "" && os.Getenv(storageFileName) == "" {
		return storage.NewMemoryRepo(), nil
	}

	return storage.NewFileRepo(config.filename)
}

func setBaseURL() {
	if config.baseURL != "" {
		return
	}

	baseURL, ok := os.LookupEnv(baseKey)
	if !ok {
		baseURL = "http://localhost:8080"
	}

	config.baseURL = baseURL
}

func setServerAddress() {
	if config.addr != "" {
		return
	}

	addr, ok := os.LookupEnv(addrKey)
	if !ok {
		addr = "localhost:8080"
	}

	config.addr = addr
}

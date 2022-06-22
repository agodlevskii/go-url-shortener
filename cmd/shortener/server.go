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
	addr     string `env:"SERVER_ADDRESS"`
	baseURL  string `env:"BASE_URL"`
	filename string `env:"FILE_STORAGE_PATH"`
}

func main() {
	flag.Parse()
	repo, err := getRepo()
	if err != nil {
		log.Fatal(err)
	}

	r := handlers.NewShortenerRouter(repo, config.baseURL)

	err = http.ListenAndServe(config.addr, r)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.StringVar(&config.addr, "a", config.addr, "The application server address")
	flag.StringVar(&config.baseURL, "b", config.baseURL, "The application server port")
	flag.StringVar(&config.filename, "f", config.filename, "The file storage name")

	setBaseURL()
	setFilename()
	setServerAddress()
}

func getRepo() (storage.Storager, error) {
	if config.filename == "" {
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

func setFilename() {
	if config.filename != "" {
		return
	}
	config.filename = os.Getenv(storageFileName)
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

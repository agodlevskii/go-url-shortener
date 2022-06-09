package main

import (
	log "github.com/sirupsen/logrus"
	"go-url-shortener/internal/handlers"
	"go-url-shortener/internal/storage"
	"net/http"
	"os"
)

const addrKey = "SERVER_ADDRESS"

func main() {
	db := storage.NewMemoryRepo()
	r := handlers.NewShortenerRouter(db)
	addr, err := getServerAddress()
	if err != nil {
		log.Error(err)
	}

	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Error(err)
	}
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

package handlers

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"go-url-shortener/configs"
	"go-url-shortener/internal/generators"
	"go-url-shortener/internal/storage"
	"go-url-shortener/internal/validators"
	"io"
	"net/http"
)

func ShortenURL(db storage.MemoRepo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil || len(b) == 0 {
			http.Error(w, "The original URL is missing. Please attach it to the request body.", http.StatusBadRequest)
			return
		}

		uri := string(b)
		if !validators.IsURLStringValid(uri) {
			http.Error(w, "You provided an incorrect URL.", http.StatusBadRequest)
			return
		}

		id, err := generateID(db, 7)
		if err != nil {
			log.Error(err)
			http.Error(w, "Couldn't generate the short URL. Please try again later.", http.StatusInternalServerError)
			return
		}

		if err = db.Add(id, uri); err != nil {
			http.Error(w, "Couldn't generate the short URL. Please try again later.", http.StatusInternalServerError)
			return
		}

		res := "http://" + configs.Host + ":" + configs.Port + "/" + id

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(res))
		if err != nil {
			log.Error(err)
		}
	}
}

func generateID(db storage.MemoRepo, size int) (string, error) {
	if size == 0 {
		size = 7
	}

	id := generators.GenerateString(size)

	for step := 1; step < 10; step++ {
		if !db.Has(id) {
			return id, nil
		}

		id = generators.GenerateString(7)
	}

	return "", errors.New("couldn't generate ID")
}

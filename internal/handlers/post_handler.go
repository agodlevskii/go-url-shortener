package handlers

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"go-url-shortener/internal/generators"
	"go-url-shortener/internal/storage"
	"go-url-shortener/internal/validators"
	"io"
	"net/http"
	"os"
)

const baseKey = "BASE_URL"

type PostRequest struct {
	URL string `json:"url"`
}

type PostResponse struct {
	Result string `json:"result"`
}

func APIPostHandler(db storage.MemoRepo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "You provided an incorrect URL request.", http.StatusBadRequest)
			log.Error(err)
			return
		}

		uri := req.URL
		if !validators.IsURLStringValid(uri) {
			http.Error(w, "You provided an incorrect URL request.", http.StatusBadRequest)
			return
		}

		shortURI, err := shortenURL(db, uri)
		if err != nil {
			log.Error(err)
			http.Error(w, "Couldn't generate the short URL. Please try again later.", http.StatusInternalServerError)
			return
		}

		res := PostResponse{Result: shortURI}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)

		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Error(err)
		}
	}
}

func WebPostHandler(db storage.MemoRepo) func(w http.ResponseWriter, r *http.Request) {
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

		res, err := shortenURL(db, uri)
		if err != nil {
			log.Error(err)
			http.Error(w, "Couldn't generate the short URL. Please try again later.", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(res))
		if err != nil {
			log.Error(err)
		}
	}
}

func shortenURL(db storage.MemoRepo, uri string) (string, error) {
	if !validators.IsURLStringValid(uri) {
		return "", errors.New("you provided an incorrect URL")
	}

	id, err := generators.GenerateID(db, 7)
	if err != nil {
		return "", err
	}
	if err = db.Add(id, uri); err != nil {
		return "", err
	}

	baseURL, err := getBaseURL()
	if err != nil {
		log.Error(err)
	}

	return baseURL + "/" + id, nil
}

func getBaseURL() (string, error) {
	var err error
	addr, ok := os.LookupEnv(baseKey)

	if !ok {
		addr = "http://localhost:8080"
		err = os.Setenv(baseKey, addr)
	}

	return addr, err
}
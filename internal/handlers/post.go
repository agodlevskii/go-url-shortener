package handlers

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"go-url-shortener/internal/generators"
	"go-url-shortener/internal/middlewares"
	"go-url-shortener/internal/storage"
	"go-url-shortener/internal/validators"
	"io"
	"net/http"
)

type PostRequest struct {
	URL string `json:"url"`
}

type PostResponse struct {
	Result string `json:"result"`
}

type UserLink struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

func APIPostHandler(db storage.Storager, baseURL string) func(w http.ResponseWriter, r *http.Request) {
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

		userID, err := middlewares.GetUserID(r)
		if err != nil {
			http.Error(w, "Cannot identify a user", http.StatusInternalServerError)
			return
		}

		shortURI, err := shortenURL(db, userID, uri, baseURL)
		if err != nil {
			log.Error(err)
			http.Error(w, "Couldn't generate the short URL. Please try again later.", http.StatusInternalServerError)
			return
		}

		res := PostResponse{Result: shortURI}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Error(err)
		}
	}
}

func WebPostHandler(db storage.Storager, baseURL string) func(w http.ResponseWriter, r *http.Request) {
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

		userID, err := middlewares.GetUserID(r)
		if err != nil {
			http.Error(w, "Cannot identify a user", http.StatusInternalServerError)
			return
		}

		res, err := shortenURL(db, userID, uri, baseURL)
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

func shortenURL(db storage.Storager, userID, uri, baseURL string) (string, error) {
	if !validators.IsURLStringValid(uri) {
		return "", errors.New("you provided an incorrect URL")
	}

	id, err := generators.GenerateID(db, userID, 7)
	if err != nil {
		return "", err
	}
	if err = db.Add(userID, id, uri); err != nil {
		return "", err
	}

	url := baseURL + "/" + id
	return url, nil
}
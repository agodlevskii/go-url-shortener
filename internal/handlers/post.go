package handlers

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"go-url-shortener/internal/apperrors"
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
			apperrors.HandleHTTPError(w, apperrors.NewError(apperrors.URLFormat, err), http.StatusBadRequest)
			return
		}

		uri := req.URL
		if !validators.IsURLStringValid(uri) {
			apperrors.HandleURLError(w)
			return
		}

		userID, err := middlewares.GetUserID(r)
		if err != nil {
			apperrors.HandleUserError(w)
			return
		}

		shortURI, chg, err := shortenURL(db, userID, uri, baseURL)
		if err != nil {
			apperrors.HandleHTTPError(w, apperrors.EmptyError(), http.StatusInternalServerError)
			return
		}

		res := PostResponse{Result: shortURI}
		w.Header().Set("Content-Type", "application/json")
		if chg {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}

		log.Error(json.NewEncoder(w).Encode(res))
	}
}

func WebPostHandler(db storage.Storager, baseURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil || len(b) == 0 {
			apperrors.HandleURLError(w)
			return
		}

		uri := string(b)
		if !validators.IsURLStringValid(uri) {
			apperrors.HandleURLError(w)
			return
		}

		userID, err := middlewares.GetUserID(r)
		if err != nil {
			apperrors.HandleUserError(w)
			return
		}

		res, chg, err := shortenURL(db, userID, uri, baseURL)
		if err != nil {
			apperrors.HandleHTTPError(w, apperrors.EmptyError(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if chg {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(res))
		if err != nil {
			log.Error(err)
		}
	}
}

func shortenURL(db storage.Storager, userID, uri, baseURL string) (string, bool, error) {
	if !validators.IsURLStringValid(uri) {
		return "", false, errors.New(apperrors.URLFormat)
	}

	id, err := generators.GenerateID(db, 7)
	if err != nil {
		return "", false, err
	}

	res, err := db.Add([]storage.ShortURL{
		{
			ID:  id,
			URL: uri,
			UID: userID,
		},
	})
	if err != nil {
		return "", false, err
	}

	url := baseURL + "/" + res[0].ID
	return url, res[0].ID != id, nil
}

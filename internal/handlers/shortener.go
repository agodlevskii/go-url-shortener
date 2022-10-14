package handlers

import (
	"encoding/json"
	"errors"
	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/generators"
	"go-url-shortener/internal/middlewares"
	"go-url-shortener/internal/storage"
	"go-url-shortener/internal/validators"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	log "github.com/sirupsen/logrus"
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

type BatchReqData struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResData struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func APIShortener(db storage.Storager, baseURL string) func(w http.ResponseWriter, r *http.Request) {
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

		if err = json.NewEncoder(w).Encode(res); err != nil {
			log.Error(err)
		}
	}
}

func WebShortener(db storage.Storager, baseURL string) func(w http.ResponseWriter, r *http.Request) {
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
		if _, err = w.Write([]byte(res)); err != nil {
			log.Error(err)
		}
	}
}

func APIBatchShortener(db storage.Storager, baseURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middlewares.GetUserID(r)
		if err != nil {
			apperrors.HandleUserError(w)
			return
		}

		var req []BatchReqData
		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			apperrors.HandleHTTPError(w, apperrors.NewError(apperrors.BatchFormat, err), http.StatusBadRequest)
			return
		}

		batch, err := getBatch(db, req, userID)
		if err != nil {
			apperrors.HandleInternalError(w)
			return
		}

		res, err := db.Add(batch)
		if err != nil {
			apperrors.HandleInternalError(w)
			return
		}

		resData := getResponseData(req, res, baseURL)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		err = json.NewEncoder(w).Encode(resData)
		if err != nil {
			apperrors.HandleInternalError(w)
		}
	}
}

func WebGetFullURL(db storage.Storager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		sURL, err := db.Get(id)
		if err != nil {
			apperrors.HandleHTTPError(w, apperrors.NewError("", err), http.StatusBadRequest)
			return
		}

		if sURL.Deleted {
			apperrors.HandleHTTPError(w, apperrors.NewError(apperrors.URLGone, nil), http.StatusGone)
			return
		}

		http.Redirect(w, r, sURL.URL, http.StatusTemporaryRedirect)
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

func getBatch(db storage.Storager, req []BatchReqData, userID string) ([]storage.ShortURL, error) {
	var batch = make([]storage.ShortURL, len(req))
	for i, data := range req {
		id, err := generators.GenerateID(db, 7)
		if err != nil {
			return nil, err
		}

		batch[i] = storage.ShortURL{
			ID:  id,
			URL: data.OriginalURL,
			UID: userID,
		}
	}

	return batch, nil
}

func getResponseData(req []BatchReqData, res []storage.ShortURL, baseURL string) []BatchResData {
	resData := make([]BatchResData, len(req))
	urlToCorID := getURLToCorIDMap(req)

	for i, sURL := range res {
		resData[i] = BatchResData{
			CorrelationID: urlToCorID[sURL.URL],
			ShortURL:      baseURL + "/" + sURL.ID,
		}
	}

	return resData
}

func getURLToCorIDMap(req []BatchReqData) map[string]string {
	res := make(map[string]string, len(req))
	for _, data := range req {
		res[data.OriginalURL] = data.CorrelationID
	}
	return res
}

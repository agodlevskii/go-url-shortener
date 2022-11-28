package handlers

import (
	"encoding/json"
	"go-url-shortener/internal/services"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"

	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/middlewares"
	"go-url-shortener/internal/storage"
)

// PostRequest describes the body for a single URL shorten request coming from API.
type PostRequest struct {
	URL string `json:"url"`
}

// PostResponse describes the response of a single URL shorten request coming from API.
type PostResponse struct {
	Result string `json:"result"`
}

// APIShortener handles the URL shortener request through API.
// The handler validates the request body to be a non-empty string of the valid format.
// It generates the shortened version and stores it in storage.ShortURL format.
func APIShortener(db storage.Storager, cfg APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middlewares.GetUserID(cfg, r)
		if err != nil {
			apperrors.HandleUserError(w)
			return
		}

		var req PostRequest
		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			apperrors.HandleHTTPError(w, apperrors.NewError(apperrors.URLFormat, err), http.StatusBadRequest)
			return
		}

		uri := req.URL
		shortURL, urlChanged, sErr := services.GetShortURL(r.Context(), db, uri, userID, cfg.GetBaseURL())
		if sErr != nil {
			handleShortenerError(w, sErr)
			return
		}

		res := PostResponse{Result: shortURL}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(getShortURLStatus(urlChanged))
		if err = json.NewEncoder(w).Encode(res); err != nil {
			log.Error(err)
		}
	}
}

// WebShortener handles the URL shortener request.
// The handler validates the request body to be a non-empty string of the valid format.
// It generates the shortened version and stores it in storage.ShortURL format.
func WebShortener(db storage.Storager, cfg APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middlewares.GetUserID(cfg, r)
		if err != nil {
			apperrors.HandleUserError(w)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil || len(b) == 0 {
			apperrors.HandleURLError(w)
			return
		}

		uri := string(b)
		shortURL, urlChanged, sErr := services.GetShortURL(r.Context(), db, uri, userID, cfg.GetBaseURL())
		if sErr != nil {
			handleShortenerError(w, sErr)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(getShortURLStatus(urlChanged))
		if _, err = w.Write([]byte(shortURL)); err != nil {
			log.Error(err)
		}
	}
}

// APIBatchShortener handles the batch URL shortener request through API.
// The handler validates the request body to match the BatchOriginalData format.
// For each provided URL, the handler generates the shortened version and stores it in storage.ShortURL format.
func APIBatchShortener(db storage.Storager, cfg APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middlewares.GetUserID(cfg, r)
		if err != nil {
			apperrors.HandleUserError(w)
			return
		}

		var data []services.BatchOriginalData
		if err = json.NewDecoder(r.Body).Decode(&data); err != nil || len(data) == 0 {
			apperrors.HandleHTTPError(w, apperrors.NewError(apperrors.BatchFormat, err), http.StatusBadRequest)
			return
		}

		resData, err := services.GetShortURLsFromBatch(r.Context(), db, data, userID, cfg.GetBaseURL())
		if err != nil {
			apperrors.HandleInternalError(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		if err = json.NewEncoder(w).Encode(resData); err != nil {
			apperrors.HandleInternalError(w)
		}
	}
}

// WebGetFullURL handles the URL redirect request.
// The handler checks if the provided shortened URL exists and not marked as deleted.
// If the validation passes, the application redirects the user to the original URL location.
func WebGetFullURL(db storage.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		sURL, err := services.GetFullURL(r.Context(), db, id)
		if err != nil {
			status := http.StatusBadRequest
			if err.Error() == apperrors.URLGone {
				status = http.StatusGone
			}

			apperrors.HandleHTTPError(w, err, status)
			return
		}

		http.Redirect(w, r, sURL.URL, http.StatusTemporaryRedirect)
	}
}

func getShortURLStatus(urlChanged bool) int {
	if urlChanged {
		return http.StatusConflict
	}
	return http.StatusCreated
}

func handleShortenerError(w http.ResponseWriter, err *apperrors.AppError) {
	if err.Error() == apperrors.URLFormat {
		apperrors.HandleURLError(w)
	} else {
		apperrors.HandleHTTPError(w, apperrors.EmptyError(), http.StatusInternalServerError)
	}
}

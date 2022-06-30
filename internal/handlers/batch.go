package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go-url-shortener/internal/generators"
	"go-url-shortener/internal/middlewares"
	"go-url-shortener/internal/storage"
	"net/http"
)

type BatchReqData struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResData struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func Batch(db storage.Storager, baseURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middlewares.GetUserID(r)
		if err != nil {
			log.Error(err)
			http.Error(w, "couldn't identify the user", http.StatusInternalServerError)
			return
		}

		var req []BatchReqData
		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			fireIncorrectRequest(w, err)
			return
		}

		var batch = make(map[string]string, len(req))
		var resData = make([]BatchResData, 0)
		for _, data := range req {
			id, err := generators.GenerateID(db, 7)
			if err != nil {
				fireInternalError(w, err)
				return
			}

			batch[id] = data.OriginalURL
			resData = append(resData, BatchResData{
				CorrelationID: data.CorrelationID,
				ShortURL:      baseURL + "/" + id,
			})
		}
		if err != nil {
			fireIncorrectRequest(w, err)
			return
		}

		err = db.AddAll(userID, batch)
		if err != nil {
			fireInternalError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		err = json.NewEncoder(w).Encode(resData)
		if err != nil {
			fireInternalError(w, err)
		}
	}
}

func fireIncorrectRequest(w http.ResponseWriter, err error) {
	http.Error(w, "You provided an incorrect batch request.", http.StatusBadRequest)
	log.Error(err)
}

func fireInternalError(w http.ResponseWriter, err error) {
	http.Error(w, "Couldn't save the passed data.", http.StatusInternalServerError)
	log.Error(err)
}

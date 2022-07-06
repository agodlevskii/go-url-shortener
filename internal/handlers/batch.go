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

		batch, err := getBatch(db, req, userID)
		if err != nil {
			fireInternalError(w, err)
			return
		}

		res, err := db.Add(batch)
		if err != nil {
			fireInternalError(w, err)
			return
		}

		resData := getResponseData(req, res, baseURL)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		err = json.NewEncoder(w).Encode(resData)
		if err != nil {
			fireInternalError(w, err)
		}
	}
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

func fireIncorrectRequest(w http.ResponseWriter, err error) {
	http.Error(w, "You provided an incorrect batch request.", http.StatusBadRequest)
	log.Error(err)
}

func fireInternalError(w http.ResponseWriter, err error) {
	http.Error(w, "Couldn't save the passed data.", http.StatusInternalServerError)
	log.Error(err)
}

package services

import (
	"context"
	"go-url-shortener/internal/generators"
	"go-url-shortener/internal/storage"
)

// BatchOriginalData describes the body for a batch URL shorten request.
// Each entity of a batch request must have a correlation ID to identify the shortened versions in the response.
// The response structure is defined in BatchTransformedData.
type BatchOriginalData struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchShortenedData describes the response of a batch URL shorten request.
// Each entity of a batch response has a correlation ID to identify the shortened versions from the request.
// The request structure is defined in BatchOriginalData.
type BatchShortenedData struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func GetShortURLsFromBatch(ctx context.Context, db storage.Storager, data []BatchOriginalData,
	uid, baseURL string) ([]BatchShortenedData, error) {
	batch, err := shortenBatchData(ctx, db, data, uid)
	if err != nil {
		return []BatchShortenedData{}, err
	}

	res, err := db.Add(ctx, batch)
	if err != nil {
		return []BatchShortenedData{}, err
	}

	return getShortenedData(data, res, baseURL), nil
}

// shortenBatchData provides the short version of each URL provided in a batch request.
// The function checks for the newly generated ID not to be associated with the existing DB entry.
func shortenBatchData(ctx context.Context, db storage.Storager,
	req []BatchOriginalData, userID string) ([]storage.ShortURL, error) {
	batch := make([]storage.ShortURL, len(req))
	for i, data := range req {
		id, err := generators.GenerateID(ctx, db, 7)
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

// getResponseData transforms the batch request into the batch response.
// Each original URL has its own ID by this moment; the function only combines the existing data.
func getShortenedData(req []BatchOriginalData, res []storage.ShortURL,
	baseURL string) []BatchShortenedData {
	resData := make([]BatchShortenedData, len(req))
	urlToCorID := getURLToCorrelationIDMap(req)

	for i, sURL := range res {
		resData[i] = BatchShortenedData{
			CorrelationID: urlToCorID[sURL.URL],
			ShortURL:      baseURL + "/" + sURL.ID,
		}
	}

	return resData
}

// getURLToCorrelationIDMap transforms the batch request into a map.
// The resulting map contains the original URL as key and the correlation ID as value.
// It is required to combine the shortened URL with associated correlation IDs.
func getURLToCorrelationIDMap(req []BatchOriginalData) map[string]string {
	res := make(map[string]string, len(req))
	for _, data := range req {
		res[data.OriginalURL] = data.CorrelationID
	}
	return res
}

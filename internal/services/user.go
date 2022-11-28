package services

import (
	"context"
	"go-url-shortener/internal/storage"
)

// UserURL describes the response for the list of all user's links.
// Each entity includes both original and shortened URLs.
type UserURL struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

// GetUserURLs recovers the user-associated links from the repository.
// In case if there are no links to return, the function returns nil instead of the empty slice.
func GetUserURLs(ctx context.Context, db storage.Storager, userID, baseURL string) ([]UserURL, error) {
	urls, err := db.GetAll(ctx, userID, false)
	if err != nil {
		return []UserURL{}, err
	}

	if len(urls) == 0 {
		return nil, nil
	}

	userURLs := make([]UserURL, 0)
	for _, url := range urls {
		userURLs = append(userURLs, UserURL{
			Short:    baseURL + "/" + url.ID,
			Original: url.URL,
		})
	}
	return userURLs, nil
}

// DeleteUserURLs deletes the user-associated links from the repository.
// The listed entities remain in the repository, but each of them gets their deletion flag set to true.
func DeleteUserURLs(ctx context.Context, db storage.Storager, userID string, ids []string) error {
	batch := make([]storage.ShortURL, 0, len(ids))
	for _, v := range ids {
		batch = append(batch, storage.ShortURL{
			ID:  v,
			UID: userID,
		})
	}

	return db.Delete(ctx, batch)
}

package services

import (
	"context"
	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/generators"
	"go-url-shortener/internal/storage"
	"go-url-shortener/internal/validators"
)

func GetFullURL(ctx context.Context, db storage.Storager, id string) (storage.ShortURL, *apperrors.AppError) {
	sURL, err := db.Get(ctx, id)
	if err != nil {
		return sURL, apperrors.NewError("", err)
	}

	if sURL.Deleted {
		return sURL, apperrors.NewError(apperrors.URLGone, nil)
	}

	return sURL, nil
}

func GetShortURL(ctx context.Context, db storage.Storager,
	uri, uid, baseURL string) (string, bool, *apperrors.AppError) {
	if !validators.IsURLStringValid(uri) {
		return "", false, apperrors.NewError(apperrors.URLFormat, nil)
	}

	id, err := generators.GenerateID(ctx, db, 7)
	if err != nil {
		return "", false, apperrors.NewError("", err)
	}

	newURLs, err := db.Add(ctx, []storage.ShortURL{
		{
			ID:  id,
			URL: uri,
			UID: uid,
		},
	})
	if err != nil {
		return "", false, apperrors.NewError("", err)
	}

	newURL := baseURL + "/" + newURLs[0].ID
	return newURL, newURLs[0].ID != id, nil
}

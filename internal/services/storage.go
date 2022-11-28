package services

import (
	"context"
	"go-url-shortener/internal/storage"
)

func Ping(ctx context.Context, db storage.Storager) bool {
	return db.Ping(ctx)
}

package services

import (
	"context"
	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/storage"
	"net"
)

type Stats struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}

func GetStatistics(ctx context.Context, db storage.Storager, ipStr, subnet string) (Stats, *apperrors.AppError) {
	var (
		network *net.IPNet
		err     error
	)

	if _, network, err = net.ParseCIDR(subnet); err != nil {
		return Stats{}, apperrors.NewError("", err)
	}

	if ip := net.ParseIP(ipStr); network.Contains(ip) {
		return collectStats(ctx, db)
	}

	return Stats{}, apperrors.NewError(apperrors.UntrustedIP, nil)
}

func collectStats(ctx context.Context, db storage.Storager) (Stats, *apperrors.AppError) {
	urls, err := db.GetAll(ctx, "", true)
	if err != nil {
		return Stats{}, apperrors.NewError("", err)
	}

	users := make(map[string]bool)
	for _, url := range urls {
		if _, ok := users[url.UID]; !ok {
			users[url.UID] = true
		}
	}

	return Stats{
		URLs:  len(urls),
		Users: len(users),
	}, nil
}

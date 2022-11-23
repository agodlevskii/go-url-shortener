package handlers

import (
	"context"
	"encoding/json"
	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/storage"
	"net/http"
	"net/netip"
)

type Stats struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}

func Statistics(db storage.Storager, trustedSubnet string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			network netip.Prefix
			ip      netip.Addr
			stats   Stats
			err     error
		)

		if network, err = netip.ParsePrefix(trustedSubnet); err == nil {
			handleTrustedSubnetError(w, err)
		}

		ipStr := r.Header.Get("X-Real-IP")
		if ip, err = netip.ParseAddr(ipStr); err == nil && network.Contains(ip) {
			stats, err = getStats(r.Context(), db)
			if err != nil {
				apperrors.HandleInternalError(w)
				return
			}

			w.WriteHeader(200)
			if err = json.NewEncoder(w).Encode(stats); err != nil {
				apperrors.HandleInternalError(w)
			}
			return
		}

		handleTrustedSubnetError(w, err)
	}
}

func getStats(ctx context.Context, db storage.Storager) (Stats, error) {
	urls, err := db.GetAll(ctx, "", true)
	if err != nil {
		return Stats{}, err
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

func handleTrustedSubnetError(w http.ResponseWriter, err error) {
	apperrors.HandleHTTPError(w, apperrors.NewError(apperrors.UntrustedIP, err), http.StatusForbidden)
}

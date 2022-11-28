package handlers

import (
	"encoding/json"
	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/services"
	"go-url-shortener/internal/storage"
	"net/http"
)

func Statistics(db storage.Storager, trustedSubnet string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ipStr := r.Header.Get("X-Real-IP")
		stats, err := services.GetStatistics(r.Context(), db, ipStr, trustedSubnet)
		if err != nil {
			apperrors.HandleHTTPError(w, err, getStatsErrorStatus(err))
			return
		}

		w.WriteHeader(200)
		if wErr := json.NewEncoder(w).Encode(stats); wErr != nil {
			apperrors.HandleHTTPError(w, apperrors.NewError("", wErr), http.StatusInternalServerError)
		}
	}
}

func getStatsErrorStatus(err error) int {
	if err.Error() == apperrors.UntrustedIP {
		return http.StatusForbidden
	}
	return http.StatusInternalServerError
}

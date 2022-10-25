package middlewares

import (
	"compress/gzip"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/respwriters"
)

// Compress provides a gzip-based encryption for the response.
// If request has valid headers, Compress replaces default http.ResponseWriter with the custom respwriters.GzipWriter.
// Otherwise, the default writer stays in charge.
func Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enc := r.Header.Get("Accept-Encoding")
		size := r.Header.Get("Content-Length")
		sizeInt, err := strconv.Atoi(size)
		if err != nil {
			sizeInt = 0
		}

		if !strings.Contains(enc, "gzip") || sizeInt < 1400 {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			apperrors.HandleHTTPError(w, apperrors.NewError("", err), http.StatusInternalServerError)
			return
		}
		defer func(gz *gzip.Writer) {
			if cErr := gz.Close(); cErr != nil {
				log.Error(cErr)
			}
		}(gz)

		gzWriter := respwriters.GzipWriter{
			ResponseWriter: w,
			Writer:         gz,
		}

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzWriter, r)
	})
}

// Decompress decrypts gzipped request.
// If request has valid header, Decompress replaces default body reader with a gzip reader.
// Otherwise, the default body reader stays in charge.
func Decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") != "gzip" {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			apperrors.HandleHTTPError(w, apperrors.NewError("", err), http.StatusInternalServerError)
			return
		}
		defer func(gz *gzip.Reader) {
			if cErr := gz.Close(); cErr != nil {
				log.Error(cErr)
			}
		}(gz)
		r.Body = gz

		next.ServeHTTP(w, r)
	})
}

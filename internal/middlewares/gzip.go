package middlewares

import (
	"compress/gzip"
	"go-url-shortener/internal/respwriters"
	"net/http"
	"strconv"
	"strings"
)

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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		gzWriter := respwriters.GzipWriter{
			ResponseWriter: w,
			Writer:         gz,
		}

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzWriter, r)
	})
}

func Decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") != "gzip" {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()
		r.Body = gz

		next.ServeHTTP(w, r)
	})
}

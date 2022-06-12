package handlers

import (
	"compress/gzip"
	"github.com/go-chi/chi/v5"
	"go-url-shortener/internal/respwriters"
	"go-url-shortener/internal/storage"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const compressFormat = "gzip"

func NewShortenerRouter(db storage.Storager, baseURL string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(compress)

	r.Route("/", func(r chi.Router) {
		r.Get("/", GetHomePage)
		r.Post("/", WebPostHandler(db, baseURL))
		r.Get("/{id}", GetFullURL(db))

		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", APIPostHandler(db, baseURL))
		})

		r.NotFound(func(writer http.ResponseWriter, request *http.Request) {
			http.Error(writer, "This HTTP method is not allowed.", http.StatusMethodNotAllowed)
		})
	})

	return r
}

func compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enc := w.Header().Get("Accept-Encoding")
		size := w.Header().Get("Content-Length")
		sizeInt, err := strconv.Atoi(size)
		if err != nil {
			sizeInt = 0
		}

		if !strings.Contains(enc, compressFormat) || sizeInt < 1400 {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		gzWriter := respwriters.GzipWriter{
			ResponseWriter: w,
			Writer:         gz,
		}

		w.Header().Set("Content-Encoding", compressFormat)
		next.ServeHTTP(gzWriter, r)
	})
}

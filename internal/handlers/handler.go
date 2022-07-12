package handlers

import (
	"github.com/go-chi/chi/v5"
	"go-url-shortener/internal/middlewares"
	"go-url-shortener/internal/storage"
	"net/http"
)

func NewShortenerRouter(db storage.Storager, baseURL string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.Authorize, middlewares.Compress, middlewares.Decompress)

	r.Route("/", func(r chi.Router) {
		r.Get("/", GetHomePage)
		r.Post("/", WebPostHandler(db, baseURL))
		r.Get("/{id}", GetFullURL(db))
		r.Get("/ping", Ping(db))

		r.Route("/api", func(r chi.Router) {
			r.Route("/shorten", func(r chi.Router) {
				r.Post("/", APIPostHandler(db, baseURL))
				r.Post("/batch", Batch(db, baseURL))
			})

			r.Route("/user", func(r chi.Router) {
				r.Route("/urls", func(r chi.Router) {
					r.Get("/", UserURLsHandler(db, baseURL))
					r.Delete("/", DeleteShortURLs(db))
				})
			})
		})

		r.NotFound(func(writer http.ResponseWriter, request *http.Request) {
			http.Error(writer, "This HTTP method is not allowed.", http.StatusMethodNotAllowed)
		})
	})

	return r
}

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
		r.Get("/ping", Ping())

		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", APIPostHandler(db, baseURL))

			r.Route("/user", func(r chi.Router) {
				r.Get("/urls", UserURLsHandler(db, baseURL))
			})
		})

		r.NotFound(func(writer http.ResponseWriter, request *http.Request) {
			http.Error(writer, "This HTTP method is not allowed.", http.StatusMethodNotAllowed)
		})
	})

	return r
}

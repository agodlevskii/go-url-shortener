package handlers

import (
	"github.com/go-chi/chi/v5"
	"go-url-shortener/internal/storage"
	"net/http"
)

func NewShortenerRouter(db storage.MemoRepo) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", GetHomePage)
		r.Post("/", WebPostHandler(db))
		r.Get("/{id}", GetFullURL(db))

		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", APIPostHandler(db))
		})

		r.NotFound(func(writer http.ResponseWriter, request *http.Request) {
			http.Error(writer, "This HTTP method is not allowed.", http.StatusMethodNotAllowed)
		})
	})

	return r
}

// Package handlers is required for registering all the application's routers.
package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/middlewares"
	"go-url-shortener/internal/storage"
)

type APIConfig interface {
	GetBaseURL() string
	GetPoolSize() int
	GetUserCookieName() string
}

// NewShortenerRouter creates a new application router with the required middleware attached.
// For the unmatched route, the handler returns Method Not Allowed response.
// The data required for the handlers' functionality is being passed to the handler or gets collected from the config.
func NewShortenerRouter(cfg APIConfig, db storage.Storager) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.Authorize(cfg), middlewares.Compress, middlewares.Decompress)
	r.Mount("/debug", middleware.Profiler())

	r.Route("/", func(r chi.Router) {
		r.Get("/", GetHomePage)
		r.Post("/", WebShortener(db, cfg))
		r.Get("/{id}", WebGetFullURL(db))
		r.Get("/ping", Ping(db))

		r.Route("/api", func(r chi.Router) {
			r.Route("/shorten", func(r chi.Router) {
				r.Post("/", APIShortener(db, cfg))
				r.Post("/batch", APIBatchShortener(db, cfg))
			})

			r.Route("/user", func(r chi.Router) {
				r.Route("/urls", func(r chi.Router) {
					r.Get("/", GetUserLinks(db, cfg))
					r.Delete("/", DeleteUserLinks(db, cfg))
				})
			})
		})

		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			apperrors.HandleHTTPError(w, apperrors.EmptyError(), http.StatusMethodNotAllowed)
		})
	})

	return r
}

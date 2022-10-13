package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/config"
	"go-url-shortener/internal/middlewares"
	"go-url-shortener/internal/storage"
	"net/http"
)

func NewShortenerRouter(db storage.Storager) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.Authorize, middlewares.Compress, middlewares.Decompress)
	r.Mount("/debug", middleware.Profiler())
	cfg := config.GetConfig()

	r.Route("/", func(r chi.Router) {
		r.Get("/", GetHomePage(cfg.Templates))
		r.Post("/", WebPostHandler(db, cfg.BaseURL))
		r.Get("/{id}", GetFullURL(db))
		r.Get("/ping", Ping(db))

		r.Route("/api", func(r chi.Router) {
			r.Route("/shorten", func(r chi.Router) {
				r.Post("/", APIPostHandler(db, cfg.BaseURL))
				r.Post("/batch", Batch(db, cfg.BaseURL))
			})

			r.Route("/user", func(r chi.Router) {
				r.Route("/urls", func(r chi.Router) {
					r.Get("/", UserURLsHandler(db, cfg.BaseURL))
					r.Delete("/", DeleteShortURLs(db, cfg.Pool))
				})
			})
		})

		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			apperrors.HandleHTTPError(w, apperrors.EmptyError(), http.StatusMethodNotAllowed)
		})
	})

	return r
}

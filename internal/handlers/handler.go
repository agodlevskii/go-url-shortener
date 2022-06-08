package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"go-url-shortener/configs"
	"go-url-shortener/internal/generators"
	"go-url-shortener/internal/storage"
	"go-url-shortener/internal/validators"
	"html/template"
	"io"
	"net/http"
)

var db = storage.NewMemoryRepo()

func NewShortenerRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", ShortenURL)
		r.Get("/", GetHomePage)
		r.Get("/{id}", GetFullURL)

		r.NotFound(func(writer http.ResponseWriter, request *http.Request) {
			http.Error(writer, "This HTTP method is not allowed.", http.StatusMethodNotAllowed)
		})
	})

	return r
}

func GetHomePage(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Error(err)
		http.Error(w, "Something went wrong. Please, try again later.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func GetFullURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	url, err := storage.GetURLFromStorage(db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		http.Error(w, "The original URL is missing. Please attach it to the request body.", http.StatusBadRequest)
		return
	}

	uri := string(b)
	if !validators.IsURLStringValid(uri) {
		http.Error(w, "You provided an incorrect URL.", http.StatusBadRequest)
		return
	}

	id, err := generateID()
	if err != nil {
		log.Error(err)
		http.Error(w, "Couldn't generate the short URL. Please try again later.", http.StatusInternalServerError)
		return
	}

	storage.AddURLToStorage(db, id, uri)
	res := "http://" + configs.Host + ":" + configs.Port + "/" + id

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res))
}

func generateID() (string, error) {
	id := generators.GenerateString(7)

	for step := 1; step < 10; step++ {
		if !db.Has(id) {
			return id, nil
		}

		id = generators.GenerateString(7)
	}

	return "", errors.New("couldn't generate ID")
}

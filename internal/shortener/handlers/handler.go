package handlers

import (
	"github.com/go-chi/chi/v5"
	"go-url-shortener/internal/shortener/storage"
	"go-url-shortener/internal/shortener/utils"
	"io"
	"net/http"
)

var db = storage.NewMemoryRepo(nil)

var index = `<html>
    <head>
    	<title>Go URL Shortener</title>
    </head>
    <body>
		<header>
        	<h1>Go URL Shortener</h1>
		</header>

		<main>
			<h2>How to use the application<h2>
			<ul>
				<li>To shorten the URL: send the POST request to this route and send the initial URL as the request body.</li>
				<li>To get the shortened URL: send the GET request and put the ID of the URL in the query parameters.</li>
			</ul>
		</main>
    </body>
</html>`

func NewShortenerRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", ShortenUrl)
		r.Get("/{id}", GetFullUrl)
		r.Get("/", GetHomePage)

		r.NotFound(func(writer http.ResponseWriter, request *http.Request) {
			http.Error(writer, "This HTTP method is not allowed.", http.StatusMethodNotAllowed)
		})
	})

	return r
}

func GetHomePage(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(index))
}

func GetFullUrl(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	url, err := storage.GetURLFromStorage(db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte(url))
}

func ShortenUrl(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		http.Error(w, "The original URL is missing. Please attach it to the request body.", http.StatusBadRequest)
		return
	}

	uri := string(b)
	if !utils.IsURLStringValid(uri) {
		http.Error(w, "You provided an incorrect URL.", http.StatusBadRequest)
		return
	}

	id := utils.GenerateString()
	storage.AddURLToStorage(db, id, uri)
	res := "http://" + r.Host + "/" + id

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(201)
	w.Write([]byte(res))
}

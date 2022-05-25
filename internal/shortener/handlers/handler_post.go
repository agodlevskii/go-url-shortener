package handlers

import (
	"go-url-shortener/internal/shortener/storage"
	"go-url-shortener/internal/shortener/utils"
	"io"
	"net/http"
)

func ShortenerPostHandler(w http.ResponseWriter, r *http.Request) {
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

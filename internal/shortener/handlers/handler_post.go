package handlers

import (
	"go-url-shortener/internal/shortener/storage"
	"io"
	"net/http"
)

func ShortenerPostHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		http.Error(w, "The original URL is missing. Please attach it to the request body.", http.StatusBadRequest)
	}

	url := string(b)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(201)
	w.Write([]byte(storage.AddUrlToStorage(url)))
}

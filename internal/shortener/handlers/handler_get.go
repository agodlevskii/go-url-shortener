package handlers

import (
	"go-url-shortener/internal/shortener/storage"
	"go-url-shortener/internal/shortener/utils"
	"net/http"
	"path"
)

func ShortenerGetHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.IsURLValid(r.URL) {
		http.Error(w, "You provided an incorrect URL.", http.StatusBadRequest)
	}

	id := path.Base(r.URL.Path)
	if id != "" && id != "/" {
		url, err := storage.GetURLFromStorage(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte(url))
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(index))
	}
}

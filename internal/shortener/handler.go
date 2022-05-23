package shortener

import (
	"io"
	"net/http"
)

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

func ShortenerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		ShortenerPostHandler(w, r)
	case http.MethodGet:
		ShortenerGetHandler(w, r)
	default:
		http.Error(w, "HTTP request method is not supported.", http.StatusMethodNotAllowed)
	}
}

func ShortenerPostHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		http.Error(w, "The original URL is missing. Please attach it to the request body.", http.StatusBadRequest)
	}

	url := string(b)
	w.WriteHeader(201)
	w.Write([]byte(AddUrlToStorage(url)))
}

func ShortenerGetHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id != "" {
		url, err := GetUrlFromStorage(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte(url))
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(index))
	}
}

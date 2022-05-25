package handlers

import (
	"go-url-shortener/internal/shortener/storage"
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

package handlers

import (
	"go-url-shortener/internal/apperrors"
	"html/template"
	"net/http"
)

// GetHomePage handles the request for the index page.
// The map of the templates is being passed from the main handler.
// If the required template is missing from the map or malformed, the user gets an error response.
func GetHomePage(tmpl map[string]*template.Template) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		home, ok := tmpl["home"]
		if !ok {
			apperrors.HandleInternalError(w)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if err := home.Execute(w, nil); err != nil {
			apperrors.HandleInternalError(w)
		}
	}
}

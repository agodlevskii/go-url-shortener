package handlers

import (
	"go-url-shortener/internal/apperrors"
	"html/template"
	"net/http"
)

func GetHomePage(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		apperrors.HandleInternalError(w)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err = tmpl.Execute(w, nil); err != nil {
		apperrors.HandleInternalError(w)
	}
}

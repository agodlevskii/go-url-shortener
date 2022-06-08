package handlers

import (
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

func GetHomePage(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Error(err)
		http.Error(w, "Something went wrong. Please, try again later.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err = tmpl.Execute(w, nil); err != nil {
		log.Error(err)
	}
}

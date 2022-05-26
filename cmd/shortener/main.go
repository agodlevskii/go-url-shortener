package main

import (
	"fmt"
	"go-url-shortener/internal/shortener/handlers"
	"net/http"
)

func main() {
	r := handlers.NewShortenerRouter()

	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		fmt.Println(err)
	}
}

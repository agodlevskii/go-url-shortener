package main

import (
	"fmt"
	"go-url-shortener/internal/shortener/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.ShortenerHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

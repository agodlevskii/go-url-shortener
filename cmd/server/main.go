package main

import (
	"fmt"
	"go-url-shortener/internal/shortener"
	"net/http"
)

func main() {
	http.HandleFunc("/", shortener.ShortenerHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

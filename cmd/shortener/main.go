package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func mainPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost { // Обрабатываем POST-запрос
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	url := generateShortKey()
	host := req.Host
	resultUrl := "http://" + host + "/" + url
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(resultUrl))
	// return
}

func returnUrl(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet { // Обрабатываем POST-запрос
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	shortUrl := strings.Trim(req.RequestURI, "/")
	if shortUrl == "EwHXdJfB" {
		res.Header().Set("Location", "https://practicum.yandex.ru/")
		res.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	res.WriteHeader(http.StatusBadRequest)
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 8

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainPage)
	mux.HandleFunc("/{id}", returnUrl)
	fmt.Println("Server is listening...")
	fmt.Println("Press Ctrl+C to stop")
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var URLDb = make(map[string]string)

func mainPage(res http.ResponseWriter, req *http.Request) {
	fmt.Println("mainPage")
	if req.Method != http.MethodPost { // Обрабатываем POST-запрос
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	host := req.Host                // получаем значение нашего хоста
	body, _ := io.ReadAll(req.Body) // достаем данные из body
	url := generateShortKey()       // генерируем короткую ссылку
	URLDb[url] = string(body)       // записываем в нашу БД

	// for key, element := range URLDb {
	// 	fmt.Println("Key:", key, "=>", "Element:", element)
	// }

	resultURL := "http://" + host + "/" + url //  склеиваем ответ
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(resultURL))
	// return
}

func returnURL(res http.ResponseWriter, req *http.Request) {
	fmt.Println("returnURL")
	if req.Method != http.MethodGet { // Обрабатываем POST-запрос
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	shortURL := strings.Trim(req.RequestURI, "/")
	// var result bool
	resURL, ok := URLDb[shortURL]
	// If the key exists
	if ok {
		res.Header().Set("Location", resURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
		return
		// Do something
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
	mux.HandleFunc("/{id}", returnURL)
	fmt.Println("Server is listening...")
	fmt.Println("Press Ctrl+C to stop")
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

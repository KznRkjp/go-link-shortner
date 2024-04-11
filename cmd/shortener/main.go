package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

var URLDb = make(map[string]string)

func mainPage(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("mainPage")
	if req.Method != http.MethodPost { // Обрабатываем POST-запрос
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// host := req.Host                // получаем значение нашего хоста
	body, _ := io.ReadAll(req.Body) // достаем данные из body
	url := generateShortKey()       // генерируем короткую ссылку
	URLDb[url] = string(body)       // записываем в нашу БД

	resultURL := flagResURL + "/" + url //  склеиваем ответ
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(resultURL))
	// return
}

func returnURL(res http.ResponseWriter, req *http.Request) {
	// fmt.Println("returnURL")
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

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rng.Intn(len(charset))]
	}
	return string(shortKey)
}

func main() {
	parseFlags()
	r := chi.NewRouter()
	r.Post("/", mainPage)
	r.Get("/{id}", returnURL)
	fmt.Println("Server is listening @", flagRunAddr)
	fmt.Println("Press Ctrl+C to stop")
	log.Fatal(http.ListenAndServe(flagRunAddr, r))
}

package app

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/KznRkjp/go-link-shortner.git/internal/flags"
	"github.com/KznRkjp/go-link-shortner.git/internal/models"
)

var URLDb = make(map[string]string)

func GetURL(res http.ResponseWriter, req *http.Request) {
	// fmt.Println("GetURL")
	if req.Method != http.MethodPost { // Обрабатываем POST-запрос
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// host := req.Host                // получаем значение нашего хоста
	body, err := io.ReadAll(req.Body) // достаем данные из body
	if err != nil {                   // валидация
		http.Error(res, "can't read body", http.StatusBadRequest)
		return
	}
	url := generateShortKey() // генерируем короткую ссылку
	URLDb[url] = string(body) // записываем в нашу БД

	resultURL := flags.FlagResURL + "/" + url //  склеиваем ответ
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(resultURL))
	// return
}

func ReturnURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet { // Обрабатываем POST-запрос
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	shortURL := strings.Trim(req.RequestURI, "/")
	// var result bool
	resURL, ok := URLDb[shortURL]
	// If the key exists
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", resURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
	// return

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

func APIGetURL(res http.ResponseWriter, req *http.Request) {
	var reqJSON models.Request
	if req.Method != http.MethodPost { // Обрабатываем POST-запрос
		res.WriteHeader(http.StatusBadRequest)
		return

	}
	// fmt.Println(req.Body)
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&reqJSON); err != nil {
		fmt.Println("parse error")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	// fmt.Println(reqJSON.URL)
	url := generateShortKey() // генерируем короткую ссылку
	URLDb[url] = reqJSON.URL  // записываем в нашу БД

	resultURL := flags.FlagResURL + "/" + url //  склеиваем ответ
	resp := models.Response{
		Result: resultURL,
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(res)
	if err := enc.Encode(resp); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

}

package app

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/KznRkjp/go-link-shortner.git/internal/flags"
)

var URLDb = make(map[string]string)

func GetURL(res http.ResponseWriter, req *http.Request) {
	fmt.Println("GetURL")
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
	fmt.Println("ReturnURL")
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

package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/KznRkjp/go-link-shortner.git/internal/database"
	"github.com/KznRkjp/go-link-shortner.git/internal/filesio"
	"github.com/KznRkjp/go-link-shortner.git/internal/flags"
	"github.com/KznRkjp/go-link-shortner.git/internal/models"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var URLDb = make(map[string]filesio.URLRecord)

// Useless at this point.
func chekIfExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {
		fmt.Println("data file exists")

	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("file does not exist")

	} else {
		fmt.Println("Dragons be there")
	}
	return true
}

// Load data from file containg json records to our im memory DB
func LoadDB(fileName string) {
	chekIfExists(fileName)
	dat, err := os.ReadFile(fileName)
	check(err)
	newDat := strings.Split(string(dat), "\n")

	consumer, err := filesio.NewConsumer(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	for i := 0; i < len(newDat)-1; i++ {
		readEvent, err := consumer.ReadEvent()
		if err != nil {
			log.Panic(err)
		}
		URLDb[readEvent.ShortURL] = *readEvent
	}
}

func GetURL(res http.ResponseWriter, req *http.Request) {

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

	flags.ParseFlags()
	if flags.FlagDBString != "" {
		database.WriteToDB(url, string(body))

	} else if len(flags.FlagDBFilePath) > 1 {

		URLDb[url] = filesio.URLRecord{ID: uint(len(URLDb)), ShortURL: url, OriginalURL: string(body)}
		//record to file if path is not empty

		producer, err := filesio.NewProducer(flags.FlagDBFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer producer.Close()
		if err := producer.WriteEvent(&filesio.URLRecord{ID: uint(len(URLDb)), ShortURL: url, OriginalURL: string(body)}); err != nil {
			log.Fatal(err)
		}

	}

	resultURL := flags.FlagResURL + "/" + url //  склеиваем ответ
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(resultURL))

}

func ReturnURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet { // Обрабатываем POST-запрос
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	shortURL := strings.Trim(req.RequestURI, "/")
	// var result bool
	resURL, ok := URLDb[shortURL]
	fmt.Println(resURL.OriginalURL)
	// If the key exists
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", resURL.OriginalURL) //  !!!!
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
	// URLDb[url] = reqJSON.URL  // записываем в нашу БД
	URLDb[url] = filesio.URLRecord{ID: uint(len(URLDb)), ShortURL: url, OriginalURL: reqJSON.URL}

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
	//record to file if path is not empty
	if len(flags.FlagDBFilePath) > 1 {
		Producer, err := filesio.NewProducer(flags.FlagDBFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer Producer.Close()
		if err := Producer.WriteEvent(&filesio.URLRecord{ID: uint(len(URLDb)), ShortURL: url, OriginalURL: reqJSON.URL}); err != nil {
			log.Fatal(err)
		}
	}

}

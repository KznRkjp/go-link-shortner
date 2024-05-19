package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	// "math/rand"
	"net/http"
	"os"
	"strings"

	// "time"

	"github.com/KznRkjp/go-link-shortner.git/internal/database"
	"github.com/KznRkjp/go-link-shortner.git/internal/filesio"
	"github.com/KznRkjp/go-link-shortner.git/internal/flags"
	"github.com/KznRkjp/go-link-shortner.git/internal/models"
	"github.com/KznRkjp/go-link-shortner.git/internal/urlgen"
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

func saveData(body []byte) string {

	url := urlgen.GenerateShortKey()

	if flags.FlagDBString != "" {
		database.WriteToDB(url, string(body), "nil")

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
	return resultURL
}

func saveDataAPI(url string, shortURL string) string {

	if flags.FlagDBString != "" {
		database.WriteToDB(url, shortURL, "nil")

	} else if len(flags.FlagDBFilePath) > 1 {
		// URLDb[url] = reqJSON.URL  // записываем в нашу БД
		// URLDb[url] = filesio.URLRecord{ID: uint(len(URLDb)), ShortURL: url, OriginalURL: reqJSON.URL}

		URLDb[url] = filesio.URLRecord{ID: uint(len(URLDb)), ShortURL: url, OriginalURL: shortURL}
		//record to file if path is not empty

		producer, err := filesio.NewProducer(flags.FlagDBFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer producer.Close()
		if err := producer.WriteEvent(&filesio.URLRecord{ID: uint(len(URLDb)), ShortURL: url, OriginalURL: shortURL}); err != nil {
			log.Fatal(err)
		}
	}
	resultURL := flags.FlagResURL + "/" + url //  склеиваем ответ
	fmt.Println(URLDb)
	return resultURL
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

	resultURL := saveData(body)
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

	if flags.FlagDBString != "" {

		resURL, err := database.GetFromDB(shortURL)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		res.Header().Set("Location", resURL)

	} else if len(flags.FlagDBFilePath) > 1 {

		resURLStruct, ok := URLDb[shortURL]
		resURL := resURLStruct.OriginalURL

		if !ok {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		res.Header().Set("Location", resURL)
	} else {
		resURLStruct, ok := URLDb[shortURL]
		resURL := resURLStruct.OriginalURL
		if !ok {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		res.Header().Set("Location", resURL)

	}

	res.WriteHeader(http.StatusTemporaryRedirect)

}

func APIGetURL(res http.ResponseWriter, req *http.Request) {
	var reqJSON models.Request
	if req.Method != http.MethodPost { // Обрабатываем POST-запрос
		res.WriteHeader(http.StatusBadRequest)
		return

	}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&reqJSON); err != nil {
		fmt.Println("parse error")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	shortURL, err := database.CheckForDuplicates(reqJSON.URL)
	if err != nil {
		url := urlgen.GenerateShortKey() // генерируем короткую ссылку
		resultURL := saveDataAPI(url, reqJSON.URL)
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
	} else {

		resp := models.Response{
			Result: flags.FlagResURL + "/" + shortURL,
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusConflict)
		enc := json.NewEncoder(res)
		if err := enc.Encode(resp); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// [{"correlation_id":"8edc229b-b33e-42ad-b5ad-41be395532c6","original_url":"http://xwn34krdyhmz6r.com/tomwj"},{"correlation_id":"9de2b71f-1279-49c9-8081-bbcbac126334","original_url":"http://xae08jvk2j.biz/phqabbnxpiy/jlvxobs77nt"}]
func APIBatchGetURL(res http.ResponseWriter, req *http.Request) {
	var sliceReqJSON []models.BatchRequest
	// var reqJSON models.BatchRequest
	if req.Method != http.MethodPost { // Обрабатываем POST-запрос
		res.WriteHeader(http.StatusBadRequest)
		return

	}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&sliceReqJSON); err != nil {
		fmt.Println("parse error")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	for i := range sliceReqJSON {
		sliceReqJSON[i].ShortURL = urlgen.GenerateShortKey()
	}
	err := database.WriteToDBBatch(sliceReqJSON)
	if err != nil {
		fmt.Println("error")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	// for i, s := range sliceReqJSON {
	// 	fmt.Println(i, s)
	// }
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(res)
	var resp []models.BatchResponse
	for i := range sliceReqJSON {
		var newResponseRecord models.BatchResponse
		newResponseRecord.CorrelationID = sliceReqJSON[i].CorrelationID
		newResponseRecord.URL = flags.FlagResURL + "/" + sliceReqJSON[i].ShortURL
		resp = append(resp, newResponseRecord)
	}
	if err := enc.Encode(resp); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

}

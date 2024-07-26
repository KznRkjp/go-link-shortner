// Приложение позволяет соеращать ссылки
// и делать прочее
package main

import (
	"database/sql"
	"log"
	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof

	"os"

	"github.com/KznRkjp/go-link-shortner.git/internal/app"
	"github.com/KznRkjp/go-link-shortner.git/internal/database"
	"github.com/KznRkjp/go-link-shortner.git/internal/flags"
	"github.com/KznRkjp/go-link-shortner.git/internal/middleware/middlelogger"
	"github.com/KznRkjp/go-link-shortner.git/internal/router"
)

// pprof
const (
	addr = ":8082" // адрес сервера pprof
)

func main() {
	flags.ParseFlags()
	if flags.FlagDBString != "" {
		var err error
		database.DB, err = sql.Open("pgx", flags.FlagDBString) // выбор способа храненеи данных в зависимости от флага.
		if err != nil {
			log.Fatal(err)
		}
		defer database.DB.Close()
		database.CreateTable(database.DB) // создание необходимых таблиц если их нет.
	} else if len(flags.FlagDBFilePath) > 0 {
		_, err := os.OpenFile(flags.FlagDBFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		app.LoadDB(flags.FlagDBFilePath)
	}
	dd := router.Main()
	go http.ListenAndServe(addr, nil) // go рутина pprof
	// записываем в лог, что сервер запускается
	middlelogger.ServerStartLog(flags.FlagRunAddr)
	if err := http.ListenAndServe(flags.FlagRunAddr, dd); err != nil {
		// записываем в лог ошибку, если сервер не запустился
		middlelogger.ServerStartLog(err.Error())
	}

}

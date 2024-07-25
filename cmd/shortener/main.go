// Oh this is sooo bad
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

const (
	addr = ":8082" // адрес сервера
)

func main() {

	flags.ParseFlags()
	if flags.FlagDBString != "" {
		var err error
		database.DB, err = sql.Open("pgx", flags.FlagDBString)
		if err != nil {
			log.Fatal(err)
		}
		defer database.DB.Close()
		database.CreateTable(database.DB)

	} else if len(flags.FlagDBFilePath) > 0 {
		_, err := os.OpenFile(flags.FlagDBFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		app.LoadDB(flags.FlagDBFilePath)
	}
	dd := router.Main()
	go http.ListenAndServe(addr, nil)

	// записываем в лог, что сервер запускается
	middlelogger.ServerStartLog(flags.FlagRunAddr)
	// defer shortlogger.Sugar.Sync()

	if err := http.ListenAndServe(flags.FlagRunAddr, dd); err != nil {
		// записываем в лог ошибку, если сервер не запустился
		middlelogger.ServerStartLog(err.Error())
	}

}

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/KznRkjp/go-link-shortner.git/internal/app"
	"github.com/KznRkjp/go-link-shortner.git/internal/database"
	"github.com/KznRkjp/go-link-shortner.git/internal/flags"
	"github.com/KznRkjp/go-link-shortner.git/internal/middleware/middlelogger"
	"github.com/KznRkjp/go-link-shortner.git/internal/router"
)

func main() {

	flags.ParseFlags()
	if flags.FlagDBString != "" {
		database.CreateTable()

	} else if len(flags.FlagDBFilePath) > 0 {
		_, err := os.OpenFile(flags.FlagDBFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		app.LoadDB(flags.FlagDBFilePath)
	}
	dd := router.Main()

	// записываем в лог, что сервер запускается
	middlelogger.ServerStartLog(flags.FlagRunAddr)
	// defer shortlogger.Sugar.Sync()

	if err := http.ListenAndServe(flags.FlagRunAddr, dd); err != nil {
		// записываем в лог ошибку, если сервер не запустился
		middlelogger.ServerStartLog(err.Error())
	}

}

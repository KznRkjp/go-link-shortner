package main

import (
	"net/http"

	"github.com/KznRkjp/go-link-shortner.git/internal/flags"
	"github.com/KznRkjp/go-link-shortner.git/internal/middleware/middlelogger"
	"github.com/KznRkjp/go-link-shortner.git/internal/router"
)

func main() {

	flags.ParseFlags()
	dd := router.Main()
	//shortlogger.ServerStartLog(flags.FlagRunAddr)

	// fmt.Println("Server is listening @", flags.FlagRunAddr)
	// fmt.Println("Press Ctrl+C to stop")
	// log.Fatal(http.ListenAndServe(flags.FlagRunAddr, dd))
	// записываем в лог, что сервер запускается
	middlelogger.ServerStartLog(flags.FlagRunAddr)
	// defer shortlogger.Sugar.Sync()

	if err := http.ListenAndServe(flags.FlagRunAddr, dd); err != nil {
		// записываем в лог ошибку, если сервер не запустился
		middlelogger.ServerStartLog(flags.FlagRunAddr)
	}

}

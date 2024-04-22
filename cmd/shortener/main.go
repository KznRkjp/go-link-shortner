package main

import (
	"net/http"

	"github.com/KznRkjp/go-link-shortner.git/internal/flags"
	"github.com/KznRkjp/go-link-shortner.git/internal/router"

	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func main() {
	// создаём предустановленный регистратор zap
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	sugar = *logger.Sugar()
	defer logger.Sync()
	flags.ParseFlags()
	dd := router.Main()

	// fmt.Println("Server is listening @", flags.FlagRunAddr)
	// fmt.Println("Press Ctrl+C to stop")
	// log.Fatal(http.ListenAndServe(flags.FlagRunAddr, dd))
	// записываем в лог, что сервер запускается
	sugar.Infow(
		"Starting server",
		"addr", flags.FlagRunAddr,
	)
	if err := http.ListenAndServe(flags.FlagRunAddr, dd); err != nil {
		// записываем в лог ошибку, если сервер не запустился
		sugar.Fatalw(err.Error(), "event", "start server")
	}
}

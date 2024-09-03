// Приложение позволяет соеращать ссылки
// и делать прочее
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof
	"os/signal"
	"syscall"

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

var buildVersion string
var buildDate string
var buildCommit string

func main() {

	printInfo()

	flags.ParseFlags()

	//DB
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
	// через этот канал сообщим основному потоку, что соединения закрыты
	idleConnsClosed := make(chan struct{})
	// канал для перенаправления прерываний
	// поскольку нужно отловить всего одно прерывание,
	// ёмкости 1 для канала будет достаточно
	sigint := make(chan os.Signal, 1)
	// регистрируем перенаправление прерываний
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	server := &http.Server{
		Handler: dd,
	}

	profServer := &http.Server{
		Handler: dd,
	}

	//Вот тут мы стартуем в HTTPS если есть флаг
	if flags.FlagHTTPSBool {
		profServer.Addr = ":4443"
		log.Println("pprof on port", profServer.Addr)
		go profServer.ListenAndServeTLS("server.crt", "server.key") // go рутина pprof
		// записываем в лог, что сервер запускается
		server.Addr = ":443"
		middlelogger.ServerStartLog(server.Addr)

		go func() {
			err := server.ListenAndServeTLS("server.crt", "server.key")
			if err != nil {
				log.Println(err)
			}
		}()

	} else {

		profServer.Addr = ":8081"
		log.Println("pprof on port", profServer.Addr)
		go profServer.ListenAndServe() // go рутина pprof
		// записываем в лог, что сервер запускается
		server.Addr = flags.FlagRunAddr
		middlelogger.ServerStartLog(flags.FlagRunAddr)
		go func() {
			if err := server.ListenAndServe(); err != nil {
				// записываем в лог ошибку, если сервер не запустился
				middlelogger.ServerStartLog(err.Error())
			}
		}()
	}
	go func() {
		// читаем из канала прерываний
		// поскольку нужно прочитать только одно прерывание,
		// можно обойтись без цикла
		<-sigint
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		if err := server.Shutdown(context.Background()); err != nil {
			// ошибки закрытия Listener
			log.Printf("HTTP server Shutdown: %v", err)
		}
		// сообщаем основному потоку,
		// что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
	}()
	// ждём завершения процедуры graceful shutdown
	<-idleConnsClosed
	// получили оповещение о завершении
	// здесь можно освобождать ресурсы перед выходом,
	// например закрыть соединение с базой данных,
	// закрыть открытые файлы
	fmt.Println("Server Shutdown gracefully")

}

func printInfo() {
	if buildVersion != "" {
		fmt.Println("Build version: ", buildVersion)
	} else {
		fmt.Println("Build version: N/A")
	}
	if buildDate != "" {
		fmt.Println("Build date: ", buildDate)
	} else {
		fmt.Println("Build date: N/A")
	}
	if buildCommit != "" {
		fmt.Println("Build commit: ", buildCommit)
	} else {
		fmt.Println("Build commit: N/A")
	}

}

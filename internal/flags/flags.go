package flags

import (
	"flag"
	"fmt"
	"os"
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var FlagRunAddr string
var FlagResURL string
var FlagDBFilePath string
var FlagDBString string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	// регистрируем переменную resURL
	flag.StringVar(&FlagResURL, "b", "http://localhost:8080", "result URL")
	// регистрируем переменную DBFilePath
	flag.StringVar(&FlagDBFilePath, "f", "/tmp/short-url-db.json", "Full path to DB file")
	// регистрируем переменную FlagDBString - для подлкючения к базе данных
	ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		`localhost`, `url`, `dwwq34zf!3`, `url`)
	flag.StringVar(&FlagDBString, "d", ps, "String for DB connection")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}
	if envResURL := os.Getenv("BASE_URL"); envResURL != "" {
		FlagResURL = envResURL
	}
	if envDBFilePath := os.Getenv("FILE_STORAGE_PATH"); envDBFilePath != "" {
		FlagDBFilePath = envDBFilePath
	}
	if envDBString := os.Getenv("DATABASE_DSN"); envDBString != "" {
		FlagDBString = envDBString
	}
}

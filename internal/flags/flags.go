package flags

import (
	"flag"
	"os"
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var FlagRunAddr string
var FlagResURL string
var FlagDBFilePath string

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

	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}
	if envResURL := os.Getenv("BASE_URL"); envResURL != "" {
		FlagResURL = envResURL
	}
	if envDBFilePath := os.Getenv("FILE_STORAGE_PATH"); envDBFilePath != "" {
		FlagResURL = envDBFilePath
	}
}

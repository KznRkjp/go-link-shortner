package main

import (
	"flag"
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var flagRunAddr string
var flagResURL string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	// регистрируем переменную resURL
	flag.StringVar(&flagResURL, "b", "http://localhost:8080", "result URL")

	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}

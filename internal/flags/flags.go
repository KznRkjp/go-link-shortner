package flags

import (
	"flag"
	"os"

	"github.com/KznRkjp/go-link-shortner.git/internal/config"
)

// FlagConfigPath содержит путь к файлу конфигурации
var FlagConfigPath string

// FlagRunAddr содержит адрес и порт для запуска сервера
var FlagRunAddr string

// FlagResURL содержит адрес сервера для сокращенной ссылки
var FlagResURL string

// FlagDBFilePath содержит данные для подключения к файлу с данными
var FlagDBFilePath string

// FlagDBString содержит данные для подключения к БД
var FlagDBString string

// FlagHTTPSString - при наличии запускает сервер в режиме HTTPS
var FlagHTTPSBool bool

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() {
	// регистрируем переменную resURL
	flag.StringVar(&FlagConfigPath, "c", "", "path to config file")

	flag.StringVar(&FlagConfigPath, "config", "", "path to config file")
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	// регистрируем переменную resURL
	flag.StringVar(&FlagResURL, "b", "http://localhost:8080", "result URL")
	// регистрируем переменную DBFilePath
	flag.StringVar(&FlagDBFilePath, "f", "/tmp/short-url-db.json", "Full path to DB file")
	// регистрируем переменную FlagDBString - для подлкючения к базе данных
	flag.StringVar(&FlagDBString, "d", "", "String for DB connection")
	// регистрируем переменную FlagHTTPSString
	flag.BoolVar(&FlagHTTPSBool, "s", false, "HTTPS mode")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if envConfigPath := os.Getenv("CONFIG"); envConfigPath != "" {
		FlagConfigPath = envConfigPath
	}
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
	if envHTTPSString := os.Getenv("ENABLE_HTTPS"); envHTTPSString != "" {
		FlagHTTPSBool = true
	}

	if FlagConfigPath != "" {
		configuration, err := config.OpenConfigFile(FlagConfigPath)
		if err == nil {
			if FlagRunAddr == "" {
				FlagRunAddr = configuration.ServerAddress
			}
			if FlagResURL == "" {
				FlagResURL = configuration.BaseURL
			}
			if FlagDBFilePath == "" {
				FlagDBFilePath = configuration.FileStoragePath
			}
			if FlagDBString == "" {
				FlagDBString = configuration.DatabaseDSN
			}
			if !FlagHTTPSBool && os.Getenv("DATABASE_DSN") == "" {
				FlagHTTPSBool = configuration.EnableHTTPS
			}

		}
	}
	// ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
	// `localhost`, `url`, `dwwq34zf!3`, `url`)
	// fmt.Println(ps)
}

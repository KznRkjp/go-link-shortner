package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Стуркрута для разбора файла конфигурации
type Config struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
}

// func main() {
// 	OpenConfigFile("config.json")
// }

// Разбор файла конфигурации
func OpenConfigFile(filename string) (Config, error) {
	var Conf Config
	f, err := os.Open(filename)
	if err != nil {
		return Conf, err
	}

	fileByte, _ := io.ReadAll(f)
	json.Unmarshal(fileByte, &Conf)
	fmt.Println(Conf)
	return Conf, nil
}

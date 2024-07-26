package filesio

import (
	"encoding/json"
	"os"
)

// URLRecord - структура описывющая данные URL
type URLRecord struct {
	ID          uint   `json:"uuid,string"`  //уникальный нормер (возможно)
	ShortURL    string `json:"short_url"`    // сокрвщенная ссылка
	OriginalURL string `json:"original_url"` // оргинальный URL
	DeletedFlag bool   `json:"is_deleted"`   // флаг удаления
}

// Producer - структура данных для записи
type Producer struct {
	file *os.File
	// добавляем Writer в Producer
	encoder *json.Encoder
}

// NewProducer - функция создания нового Producer для последующей записи
func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file: file,
		// создаём новый Writer
		encoder: json.NewEncoder(file),
	}, nil
}

// WriteEvent - метод для записи нового события
func (p *Producer) WriteEvent(event *URLRecord) error {
	return p.encoder.Encode(&event)
}

// Close - метод для закрытия файла
func (p *Producer) Close() error {
	// закрываем файл
	return p.file.Close()
}

// Consumer - структура для обращения к файлу данных
type Consumer struct {
	file *os.File
	// добавляем reader в Consumer
	decoder *json.Decoder
}

// NewConsumer и NewConsumer..
func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file: file,
		// создаём новый scanner
		decoder: json.NewDecoder(file),
	}, nil
}

// Consumer.ReadEvent - построчное чтение данных
func (c *Consumer) ReadEvent() (*URLRecord, error) {
	// одиночное сканирование до следующей строки
	event := &URLRecord{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

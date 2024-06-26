package filesio

import (
	"encoding/json"
	"os"
)

type URLRecord struct {
	ID          uint   `json:"uuid,string"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	DeletedFlag bool   `json:"is_deleted"`
}

type Producer struct {
	file *os.File
	// добавляем Writer в Producer
	encoder *json.Encoder
}

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

func (p *Producer) WriteEvent(event *URLRecord) error {
	return p.encoder.Encode(&event)
}

func (p *Producer) Close() error {
	// закрываем файл
	return p.file.Close()
}

type Consumer struct {
	file *os.File
	// добавляем reader в Consumer
	decoder *json.Decoder
}

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

package filesio

import (
	"testing"
)

func TestNewProducer(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *Producer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProducer(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProducer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewProducer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConsumer(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *Consumer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConsumer(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConsumer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewConsumer() = %v, want %v", got, tt.want)
			}
		})
	}
}

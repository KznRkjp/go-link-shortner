package middlelogger

import (
	"net/http"
	"reflect"
	"testing"
)

func TestServerStartLog(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Test1", args{"http://0.0.0.0:8000"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ServerStartLog(tt.args.addr)
		})
	}
}

func TestServerErrorLog(t *testing.T) {
	type args struct {
		error string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Test2", args{"Really nasty error"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ServerErrorLog(tt.args.error)
		})
	}
}

func TestWithLogging(t *testing.T) {

	type args struct {
		h http.Handler
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithLogging(tt.args.h); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithLogging() = %v, want %v", got, tt.want)
			}
		})
	}
}

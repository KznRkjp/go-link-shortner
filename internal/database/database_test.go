package database

import (
	"net/http"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestPing(t *testing.T) {
	type args struct {
		res http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Ping(tt.args.res, tt.args.req)
		})
	}
}

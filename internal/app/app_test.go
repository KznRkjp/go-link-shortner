package app

import (
	"net/http"
	"testing"
)

func TestAPIGetURL(t *testing.T) {
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
			APIGetURL(tt.args.res, tt.args.req)
		})
	}
}

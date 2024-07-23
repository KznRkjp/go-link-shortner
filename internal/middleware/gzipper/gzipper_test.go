package gzipper

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_newCompressWriter(t *testing.T) {
	type args struct {
		w http.ResponseWriter
	}
	tests := []struct {
		name string
		args args
		want *compressWriter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newCompressWriter(tt.args.w); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCompressWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}

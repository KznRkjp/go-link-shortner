package users

import (
	"net/http"
	"testing"
)

func TestBuildJWTString(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildJWTString(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildJWTString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BuildJWTString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUserUID(t *testing.T) {
	type args struct {
		tokenString string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserUID(tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccess(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Access(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Access() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Access() = %v, want %v", got, tt.want)
			}
		})
	}
}

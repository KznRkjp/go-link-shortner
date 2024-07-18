package app

import (
	"context"
	"testing"
)

func Test_saveData(t *testing.T) {
	type args struct {
		ctx  context.Context
		body []byte
		uuid string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := saveData(tt.args.ctx, tt.args.body, tt.args.uuid); got != tt.want {
				t.Errorf("saveData() = %v, want %v", got, tt.want)
			}
		})
	}
}

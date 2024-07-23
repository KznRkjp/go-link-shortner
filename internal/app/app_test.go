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
		want int
	}{
		{
			name: "Test #1",
			args: args{
				ctx:  context.Background(),
				body: []byte("http://mail.ru"),
				uuid: "1dsf24123123123",
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := saveData(tt.args.ctx, tt.args.body, tt.args.uuid); len([]rune(got)) < tt.want {
				t.Errorf("saveData() = %v, length is less than %v", got, tt.want)
			}
		})
	}
}

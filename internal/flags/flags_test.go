package flags

import "testing"

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Check if it works",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ParseFlags()
		})
	}
}

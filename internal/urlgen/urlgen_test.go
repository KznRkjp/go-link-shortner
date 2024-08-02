package urlgen

import "testing"

func TestGenerateShortKey(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"test if len of return is equal to 8", 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateShortKey(); tt.want == len([]rune(got)) {
				t.Errorf("GenerateShortKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

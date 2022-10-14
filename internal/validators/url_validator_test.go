package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsURLStringValid(t *testing.T) {
	tests := []struct {
		name   string
		rawURL string
		want   bool
	}{
		{
			name:   "Incorrect URL",
			rawURL: "URL",
			want:   false,
		},
		{
			name:   "URL without protocol",
			rawURL: "google.com",
			want:   false,
		},
		{
			name:   "Correct URL",
			rawURL: "https://google.com",
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsURLStringValid(tt.rawURL))
		})
	}
}

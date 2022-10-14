package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	if err := os.Chdir("../../"); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name          string
		want          *Config
		wantTemplates int
	}{
		{
			name: "Default config",
			want: &Config{
				Addr:    "localhost:8080",
				BaseURL: "http://localhost:8080",
				Pool:    10,
			},
			wantTemplates: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetConfig()
			assert.Equal(t, tt.want.Addr, got.Addr)
			assert.Equal(t, tt.want.BaseURL, got.BaseURL)
			assert.Equal(t, tt.want.Pool, got.Pool)
			assert.Equal(t, tt.wantTemplates, len(got.Templates))
		})
	}
}

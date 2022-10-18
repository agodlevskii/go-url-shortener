package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "Default config",
			want: &Config{
				Addr:    "localhost:8080",
				BaseURL: "http://localhost:8080",
				Pool:    10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetConfig()
			assert.Equal(t, tt.want.Addr, got.Addr)
			assert.Equal(t, tt.want.BaseURL, got.BaseURL)
			assert.Equal(t, tt.want.Pool, got.Pool)
		})
	}
}

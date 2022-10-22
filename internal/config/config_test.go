package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "Default config",
			want: &Config{PoolSize: 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New()
			assert.Equal(t, tt.want.Addr, got.Addr)
			assert.Equal(t, tt.want.BaseURL, got.BaseURL)
			assert.Equal(t, tt.want.PoolSize, got.PoolSize)
		})
	}
}

func TestWithEnv(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "Default env config",
			want: &Config{
				Addr:     "localhost:8080",
				BaseURL:  "http://localhost:8080",
				PoolSize: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(WithEnv())
			assert.Equal(t, tt.want.Addr, got.Addr)
			assert.Equal(t, tt.want.BaseURL, got.BaseURL)
			assert.Equal(t, tt.want.PoolSize, got.PoolSize)
		})
	}
}

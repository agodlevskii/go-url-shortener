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

func TestWithFlags(t *testing.T) {
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
			got := New(WithFlags())
			assert.Equal(t, tt.want.Addr, got.Addr)
			assert.Equal(t, tt.want.BaseURL, got.BaseURL)
			assert.Equal(t, tt.want.PoolSize, got.PoolSize)
		})
	}
}

func TestConfig_GetBaseURL(t *testing.T) {
	cfg := New(WithEnv())
	assert.Equal(t, "http://localhost:8080", cfg.GetBaseURL())
}

func TestConfig_GetServerAddr(t *testing.T) {
	cfg := New(WithEnv())
	assert.Equal(t, "localhost:8080", cfg.GetServerAddr())
}

func TestConfig_GetStorageFileName(t *testing.T) {
	cfg := New(WithEnv())
	assert.Equal(t, "", cfg.GetStorageFileName())
}

func TestConfig_GetDBURL(t *testing.T) {
	cfg := New(WithEnv())
	assert.Equal(t, "", cfg.GetDBURL())
}

func TestConfig_GetPoolSize(t *testing.T) {
	cfg := New(WithEnv())
	assert.Equal(t, 10, cfg.GetPoolSize())
}

func TestConfig_GetUserCookieName(t *testing.T) {
	cfg := New(WithEnv())
	assert.Equal(t, "user_id", cfg.GetUserCookieName())
}

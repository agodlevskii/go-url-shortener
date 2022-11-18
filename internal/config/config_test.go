package config

import (
	"encoding/json"
	"os"
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

func TestWithFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		fileCfg  Config
		want     *Config
	}{
		{
			name: "Default config",
			want: &Config{
				PoolSize:       10,
				UserCookieName: "user_id",
			},
		},
		{
			name:     "File config",
			filename: "test_cfg.json",
			fileCfg:  Config{Secure: true},
			want: &Config{
				Addr:           "localhost:8080",
				BaseURL:        "http://localhost:8080",
				ConfigFile:     "test_cfg.json",
				PoolSize:       10,
				Secure:         true,
				UserCookieName: "user_id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *Config

			if tt.filename != "" {
				if err := setupFileConfig(tt.filename, tt.fileCfg); err != nil {
					t.Fatal(err)
				}
				got = New(WithEnv(), WithFile())
			} else {
				got = New(WithFile())
			}

			assert.Equal(t, tt.want.Addr, got.Addr)
			assert.Equal(t, tt.want.BaseURL, got.BaseURL)
			assert.Equal(t, tt.want.ConfigFile, got.ConfigFile)
			assert.Equal(t, tt.want.DBURL, got.DBURL)
			assert.Equal(t, tt.want.Filename, got.Filename)
			assert.Equal(t, tt.want.PoolSize, got.PoolSize)
			assert.Equal(t, tt.want.Secure, got.Secure)
			assert.Equal(t, tt.want.UserCookieName, got.UserCookieName)

			if tt.filename != "" {
				if err := cleanFileConfig(tt.filename); err != nil {
					t.Fatal(err)
				}
			}
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

func TestConfig_IsSecure(t *testing.T) {
	cfg := New(WithEnv())
	assert.Equal(t, false, cfg.IsSecure())
}

func setupFileConfig(filename string, cfg Config) error {
	if err := os.Setenv("CONFIG", filename); err != nil {
		return err
	}

	var f *os.File
	var err error
	f, err = os.Create(filename)
	if err != nil {
		return err
	}
	return json.NewEncoder(f).Encode(cfg)
}

func cleanFileConfig(filename string) error {
	if err := os.Unsetenv("CONFIG"); err != nil {
		return err
	}
	return os.Remove(filename)
}

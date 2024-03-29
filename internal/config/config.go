// Package config includes the shared application configuration values.
package config

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"reflect"

	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
)

// Config describes the configuration required across the application.
// Since the configuration can be initiated via the environment flags, the struct contains the required annotation.
type Config struct {
	Addr           string `json:"server_address" env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL        string `json:"base_url" env:"BASE_URL" envDefault:"http://localhost:8080"`
	ConfigFile     string `env:"CONFIG"`
	DBURL          string `json:"database_dsn" env:"DATABASE_DSN"`
	Filename       string `json:"file_storage_path" env:"FILE_STORAGE_PATH"`
	PoolSize       int    `json:"pool_size" env:"POOL_SIZE" envDefault:"10"`
	Secure         bool   `json:"enable_https" env:"ENABLE_HTTPS"`
	UserCookieName string `json:"user_cookie" env:"USER_COOKIE" envDefault:"user_id"`
}

func New(opts ...func(*Config)) *Config {
	cfg := &Config{
		PoolSize:       10,
		UserCookieName: "user_id",
	}
	for _, o := range opts {
		o(cfg)
	}
	return cfg
}

func WithEnv() func(*Config) {
	return func(cfg *Config) {
		if err := env.Parse(cfg); err != nil {
			log.Fatal(err)
		}
	}
}

func WithFlags() func(*Config) {
	return func(cfg *Config) {
		flag.StringVar(&cfg.Addr, "a", cfg.Addr, "The application server address")
		flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "The application server port")
		flag.StringVar(&cfg.ConfigFile, "c", cfg.ConfigFile, "The application configuration file")
		flag.StringVar(&cfg.ConfigFile, "config", cfg.ConfigFile, "The application configuration file")
		flag.StringVar(&cfg.DBURL, "d", cfg.DBURL, "The DB connection URL")
		flag.BoolVar(&cfg.Secure, "s", cfg.Secure, "The HTTPS connection config")
		flag.StringVar(&cfg.Filename, "f", cfg.Filename, "The file storage name")
		flag.Parse()
	}
}

func WithFile() func(*Config) {
	return func(cfg *Config) {
		if cfg.ConfigFile == "" {
			return
		}
		fileCfg, err := getConfigFromFile(cfg.ConfigFile)
		if err != nil {
			log.Error(err)
			return
		}

		rCfg := reflect.Indirect(reflect.ValueOf(cfg))
		rFileCfg := reflect.ValueOf(fileCfg)
		for i := 0; i < rCfg.NumField(); i++ {
			rField := rCfg.Type().Field(i).Name
			rValue := rCfg.FieldByName(rField)

			if rValue.IsZero() && rValue.CanSet() {
				if fileValue := rFileCfg.FieldByName(rField); !fileValue.IsZero() {
					rValue.Set(fileValue)
				}
			}
		}
	}
}

func getConfigFromFile(filepath string) (Config, error) {
	var cfg Config

	file, err := os.OpenFile(filepath, os.O_RDONLY|os.O_CREATE, 0o777)
	if err != nil {
		return cfg, err
	}

	defer func(file *os.File) {
		if fErr := file.Close(); fErr != nil {
			log.Error(err)
		}
	}(file)

	rawCfg, err := io.ReadAll(file)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(rawCfg, &cfg)
	return cfg, err
}

func (c *Config) GetBaseURL() string {
	return c.BaseURL
}

func (c *Config) GetDBURL() string {
	return c.DBURL
}

func (c *Config) IsSecure() bool {
	return c.Secure
}

func (c *Config) GetPoolSize() int {
	return c.PoolSize
}

func (c *Config) GetServerAddr() string {
	return c.Addr
}

func (c *Config) GetStorageFileName() string {
	return c.Filename
}

func (c *Config) GetUserCookieName() string {
	return c.UserCookieName
}

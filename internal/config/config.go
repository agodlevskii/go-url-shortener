// Package config includes the shared application configuration values.
package config

import (
	"flag"
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
)

// Config describes the configuration required across the application.
// Since the configuration can be initiated via the environment flags, the struct contains the required annotation.
type Config struct {
	Addr           string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL        string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	DBURL          string `env:"DATABASE_DSN"`
	Filename       string `env:"FILE_STORAGE_PATH"`
	PoolSize       int
	UserCookieName string
}

func New(opts ...func(*Config)) *Config {
	cfg := &Config{
		UserCookieName: "user_id",
		PoolSize:       10,
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
		flag.StringVar(&cfg.DBURL, "d", cfg.DBURL, "The DB connection URL")
		flag.StringVar(&cfg.Filename, "f", cfg.Filename, "The file storage name")
		flag.Parse()
	}
}

func (c *Config) GetBaseURL() string {
	return c.BaseURL
}

func (c *Config) GetDBURL() string {
	return c.DBURL
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

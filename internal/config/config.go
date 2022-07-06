package config

import (
	"flag"
	"github.com/caarlos0/env"
)

type Config struct {
	Addr     string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL  string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	DBURL    string `env:"DATABASE_DSN"`
	Filename string `env:"FILE_STORAGE_PATH"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := initFromEnvVar(cfg)
	if err != nil {
		return nil, err
	}

	initFromFlags(cfg)
	return cfg, nil
}

func initFromEnvVar(cfg *Config) error {
	return env.Parse(cfg)
}

func initFromFlags(cfg *Config) {
	flag.StringVar(&cfg.Addr, "a", cfg.Addr, "The application server address")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "The application server port")
	flag.StringVar(&cfg.DBURL, "d", cfg.DBURL, "The DB connection URL")
	flag.StringVar(&cfg.Filename, "f", cfg.Filename, "The file storage name")
	flag.Parse()
}

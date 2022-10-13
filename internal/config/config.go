package config

import (
	"flag"
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
	"html/template"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	Addr      string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL   string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	DBURL     string `env:"DATABASE_DSN"`
	Filename  string `env:"FILE_STORAGE_PATH"`
	Templates map[string]*template.Template
	Pool      int
}

var (
	cfg  *Config
	once sync.Once
)

func GetConfig() *Config {
	once.Do(initConfig)
	return cfg
}

func initConfig() {
	cfg = &Config{
		Templates: make(map[string]*template.Template),
		Pool:      10,
	}
	if err := initTemplates(cfg); err != nil {
		log.Fatal(err)
	}

	if err := initFromEnvVar(cfg); err != nil {
		log.Fatal(err)
	}

	initFromFlags(cfg)
}

func initTemplates(cfg *Config) error {
	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, "templates/index.html")

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return err
	}

	cfg.Templates["home"] = tmpl
	return nil
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

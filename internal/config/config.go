// Package config includes the shared application configuration values.
package config

import (
	"flag"
	"html/template"
	"os"
	"path/filepath"
	"sync"

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
	Templates      map[string]*template.Template
	Pool           int
	UserCookieName string
}

var (
	cfg  *Config
	once sync.Once
)

// GetConfig returns the application configuration.
// The application uses a single instance of the config, so the configuration gets initiated only once.
func GetConfig() *Config {
	once.Do(initConfig)
	return cfg
}

// initConfig provides the initial configuration of the application.
// The initialization includes reading the templates from HDD, and processing environment variables and CLI flags.
// If the template reading and the environment processing fails, the application will exit with error.
func initConfig() {
	cfg = &Config{
		Templates:      make(map[string]*template.Template),
		UserCookieName: "user_id",
		Pool:           10,
	}
	if err := initTemplates(cfg); err != nil {
		log.Fatal(err)
	}

	if err := initFromEnvVar(cfg); err != nil {
		log.Fatal(err)
	}

	initFromFlags(cfg)
}

// initTemplates reads the HTML templates stored on the HDD.
// If the template fails to be parsed, the error will be returned.
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

// initFromEnvVar reads the configuration values from the environment variables.
// If the template fails to be parsed, the error will be returned.
func initFromEnvVar(cfg *Config) error {
	return env.Parse(cfg)
}

// initFromFlags reads the configuration values from the environment variables.
func initFromFlags(cfg *Config) {
	flag.StringVar(&cfg.Addr, "a", cfg.Addr, "The application server address")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "The application server port")
	flag.StringVar(&cfg.DBURL, "d", cfg.DBURL, "The DB connection URL")
	flag.StringVar(&cfg.Filename, "f", cfg.Filename, "The file storage name")
	flag.Parse()
}

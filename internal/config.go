package internal

import (
	"flag"
	"github.com/caarlos0/env"
)

var Config struct {
	Addr     string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL  string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	DBURL    string `env:"DATABASE_DSN" envDefault:"postgres://yand:yand@localhost:5432/practicum"`
	Filename string `env:"FILE_STORAGE_PATH"`
}

func InitConfig() error {
	err := env.Parse(&Config)
	if err != nil {
		return err
	}

	flag.Parse()
	return nil
}

func init() {
	flag.StringVar(&Config.Addr, "a", Config.Addr, "The application server address")
	flag.StringVar(&Config.BaseURL, "b", Config.BaseURL, "The application server port")
	flag.StringVar(&Config.DBURL, "d", Config.DBURL, "The DB connection URL")
	flag.StringVar(&Config.Filename, "f", Config.Filename, "The file storage name")
}

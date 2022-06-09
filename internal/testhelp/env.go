package testhelp

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func RemoveEnvVar(key string) {
	if _, ok := os.LookupEnv(key); !ok {
		err := os.Remove(key)
		if err != nil {
			log.Error(err)
		}
	}
}

func SetEnvVar(key, val string) {
	err := os.Setenv(key, val)
	if err != nil {
		log.Error(err)
	}
}

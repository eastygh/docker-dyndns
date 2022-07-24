package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	User     string
	Password string
	Zone     string
	Domains  []string
	TTL      string
}

func (conf *Config) ParseConfig(pathToConfig string) {
	file, err := os.Open(pathToConfig)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		panic(err)
	}
}

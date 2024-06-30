package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Server Server
}

type Server struct {
	Port int `env:"PORT" envDefault:"8005"`
}

// Read reads config.
func Read() (Config, error) {
	conf := Config{}
	if err := env.Parse(&conf); err != nil {
		return Config{}, fmt.Errorf("parse config from env: %w", err)
	}

	return conf, nil
}

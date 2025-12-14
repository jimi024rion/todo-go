package config

import (
	"github.com/caarlos0/env/v6"
)

// Config holds the application configuration.
type Config struct {
	Port int `env:"PORT" envDefault:"8080"`
}

// Load loads the configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

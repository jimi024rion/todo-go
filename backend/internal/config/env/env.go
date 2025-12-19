package env

import (
	"github.com/caarlos0/env/v6"
)

// Config holds the application configuration.
type Config struct {
	Port int      `env:"PORT" envDefault:"8080"`
	DB   DBConfig `envPrefix:"DB_"`
}

// DBConfig holds the database connection configuration.
type DBConfig struct {
	Host     string `env:"HOST,required"`
	User     string `env:"USER,required"`
	Password string `env:"PASSWORD,required"`
	Name     string `env:"NAME,required"`
	Port     int    `env:"PORT" envDefault:"5432"`
}

// Load loads the configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

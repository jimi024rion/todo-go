package env

import (
	envConfig "github.com/caarlos0/env/v6"
)

// Config holds the application configuration.
type Config struct {
	AppEnv       string   `env:"APP_ENV" envDefault:"locala"`
	Port         int      `env:"PORT" envDefault:"8080"`
	DB           DBConfig `envPrefix:"DB_"`
	GCPProjectID string   `env:"GCP_PROJECT_ID" envDefault:"test-project"`
}

// DBConfig holds the database connection configuration.
type DBConfig struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	User     string `env:"USER" envDefault:"postgres"`
	Password string `env:"PASSWORD" envDefault:"postgres"`
	Name     string `env:"NAME" envDefault:"todo-db"`
	Port     int    `env:"PORT" envDefault:"5432"`
}

// Load loads the configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := envConfig.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// GetConfig is a helper function to load the configuration.
func GetConfig() *Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

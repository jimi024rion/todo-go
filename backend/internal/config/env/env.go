package env

import (
	"fmt"

	envConfig "github.com/caarlos0/env/v11"
)

// Config holds the application configuration.
type Config struct {
	AppEnv       string   `env:"APP_ENV" envDefault:"local"`
	Port         int      `env:"PORT" envDefault:"8080"`
	DB           DBConfig `envPrefix:"DB_"`
	GCPProjectID string   `env:"GCP_PROJECT_ID" envDefault:"test-project"`
	LogLevel     string   `env:"LOG_LEVEL" envDefault:"debug"`
}

// DBConfig holds the database connection configuration.
type DBConfig struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	User     string `env:"USER" envDefault:"postgres"`
	Password string `env:"PASSWORD" envDefault:"postgres"`
	Name     string `env:"NAME" envDefault:"todo-db"`
	Port     int    `env:"PORT" envDefault:"5432"`
}

var Cfg Config // nolint:gochecknoglobals

// Load loads the configuration from environment variables.
func Load() error {
	if err := envConfig.Parse(&Cfg); err != nil {
		return err
	}
	if !isValidLogLevel(Cfg.LogLevel) {
		return fmt.Errorf("invalid LOG_LEVEL: %s", Cfg.LogLevel)
	}

	return nil
}

func isValidLogLevel(level string) bool {
	switch level {
	case "debug", "info", "warn", "error", "fatal", "panic":
		return true
	default:
		return false
	}
}

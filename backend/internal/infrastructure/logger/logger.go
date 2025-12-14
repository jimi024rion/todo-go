package logger

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Setup initializes the logger based on the GIN_MODE environment variable.
func Setup() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Use JSON logger in release mode, and console logger otherwise.
	if gin.Mode() == gin.ReleaseMode {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}
}

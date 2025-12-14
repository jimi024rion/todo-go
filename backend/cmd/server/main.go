package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/config"
	"github.com/jimi024rion/todo-go/backend/internal/infrastructure/logger"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http"
	"github.com/rs/zerolog/log"
)

func main() {
	// Set Gin mode based on environment (e.g., GIN_MODE)
	gin.SetMode(gin.ReleaseMode)

	// Setup logger
	logger.Setup()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Setup Router
	r := http.NewRouter()

	// Run the server
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Info().Msgf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

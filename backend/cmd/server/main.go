package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/cmd/server/di"
	"github.com/jimi024rion/todo-go/backend/internal/config/env"
	"github.com/jimi024rion/todo-go/backend/internal/config/logger"
	"github.com/jimi024rion/todo-go/backend/internal/config/trace"
	"github.com/rs/zerolog/log"
)

func main() {
	// Set Gin mode based on environment (e.g., GIN_MODE)
	gin.SetMode(gin.ReleaseMode)

	// Initialize logger
	logger.InitializeLogger()
	l := logger.NewLogger(context.Background())

	// Initialize tracer
	shutdownTracer, err := trace.InitializeTracer()
	if err != nil {
		l.FatalLog(err)
	}
	defer func() {
		if err := shutdownTracer(context.Background()); err != nil {
			l.FatalLog(err)
		}
	}()

	// Load configuration
	cfg, err := env.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Setup Server with DI
	server, cleanup, err := di.InitializeServer()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize server")
	}
	defer cleanup()

	// Run the server
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Info().Msgf("Starting server on %s", addr)
	if err := server.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

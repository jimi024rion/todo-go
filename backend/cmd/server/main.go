package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	di "github.com/jimi024rion/todo-go/backend/cmd/server/di/kessoku"
	"github.com/jimi024rion/todo-go/backend/internal/config/env"
	"github.com/jimi024rion/todo-go/backend/internal/config/logger"
	"github.com/jimi024rion/todo-go/backend/internal/config/trace"
	"github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb"
)

func main() {
	ctx := context.Background()

	// Set Gin mode based on environment (e.g., GIN_MODE)
	gin.SetMode(gin.ReleaseMode)

	// Load configuration
	err := env.Load()
	if err != nil {
		// Logger is not initialized yet, so use standard fmt for this critical error.
		panic(fmt.Sprintf("failed to load configuration: %v", err))
	}

	// Initialize logger
	logger.InitializeLogger()
	l := logger.NewLogger(ctx)

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

	// Initialize database
	db, err := rdb.NewDB(&env.Cfg)
	if err != nil {
		l.FatalLog(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			l.FatalLog(err)
		}
	}()

	// Setup Server with DI
	server := di.InitializeServer(ctx, db)

	// Run the server
	addr := fmt.Sprintf(":%d", env.Cfg.Port)
	l.InfoLog(fmt.Sprintf("Starting server on %s", addr))
	if err := server.Run(addr); err != nil {
		l.FatalLog(err)
	}
}

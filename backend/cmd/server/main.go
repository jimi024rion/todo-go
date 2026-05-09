package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	di "github.com/jimi024rion/todo-go/backend/cmd/server/di/kessoku"
	_ "github.com/jimi024rion/todo-go/backend/docs/swagger"
	"github.com/jimi024rion/todo-go/backend/internal/config/env"
	"github.com/jimi024rion/todo-go/backend/internal/config/logger"
	"github.com/jimi024rion/todo-go/backend/internal/config/trace"
	"github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb"
)

// @title           Todo API
// @version         1.0
// @description     Todo管理アプリケーションのREST API
// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        X-API-Key
// @description                 X-API-Keyヘッダーにトークンを設定してください

// @securityDefinitions.apikey  SkipAuth
// @in                          header
// @name                        X-Skip-Auth
// @description                 認証をスキップするためのヘッダー。テストや特定の内部APIで使用する。ユーザーIDを指定する。

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
	shutdownTracer, err := trace.InitializeTracer(ctx)
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
	server, err := di.InitializeServer(ctx, db)
	if err != nil {
		l.FatalLog(err)
	}

	// Run the server
	addr := fmt.Sprintf(":%d", env.Cfg.Port)
	l.InfoLog(fmt.Sprintf("Starting server on %s", addr))
	if err := server.Run(addr); err != nil {
		l.FatalLog(err)
	}
}

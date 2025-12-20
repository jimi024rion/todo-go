//go:build wireinject
// +build wireinject

package di

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/handler"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/health"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/todo"
)

// InitializeServer initializes the Gin engine with all its dependencies.
func InitializeServer() (*gin.Engine, func(), error) {
	wire.Build(
		http.NewRouter,
		handler.NewHandler,
		// health
		health.NewHandler,
		health.NewHealthCheckHandler,
		// todo
		todo.NewHandler,
		todo.NewListHandler,
		todo.NewCreateHandler,
		todo.NewGetHandler,
		todo.NewUpdateHandler,
		todo.NewDeleteHandler,
	)

	// ここはwire generateによって置き換えられる
	return nil, nil, nil
}

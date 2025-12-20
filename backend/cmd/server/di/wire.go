//go:build wireinject
// +build wireinject

package di

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/handler"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/health"
	todo_presen "github.com/jimi024rion/todo-go/backend/internal/presentation/http/todo"
	todo_usecase "github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
)

// InitializeServer initializes the Gin engine with all its dependencies.
func InitializeServer() (*gin.Engine, func(), error) {
	wire.Build(
		// --- presentation ---
		http.NewRouter,
		handler.NewHandler,
		// health
		health.NewHandler,
		health.NewHealthCheckHandler,
		// todo
		todo_presen.NewHandler,
		todo_presen.NewListHandler,
		todo_presen.NewCreateHandler,
		todo_presen.NewGetHandler,
		todo_presen.NewUpdateHandler,
		todo_presen.NewDeleteHandler,

		// --- usecase ---
		// todo
		todo_usecase.NewListUseCase,
		todo_usecase.NewCreateUseCase,
		todo_usecase.NewGetUseCase,
		todo_usecase.NewUpdateUseCase,
		todo_usecase.NewDeleteUseCase,
	)

	// ここはwire generateによって置き換えられる
	return nil, nil, nil
}

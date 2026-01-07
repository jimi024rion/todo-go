//go:build wireinject
// +build wireinject

package di

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	todorepo "github.com/jimi024rion/todo-go/backend/internal/domain/todo/repository"
	"github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb"
	todoinfra "github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb/repository/todo"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/handler"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/health"
	todopresen "github.com/jimi024rion/todo-go/backend/internal/presentation/http/todo"
	todousecase "github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
	"github.com/stephenafamo/bob"
)

// InitializeServer initializes the Gin engine with all its dependencies.
func InitializeServer(db bob.DB) (*gin.Engine, error) {
	wire.Build(
		// --- infrastructure ---
		// tx
		rdb.NewTxManager,
		// todo
		todoinfra.NewRepository,
		wire.Bind(new(todorepo.TodoRepository), new(*todoinfra.Repository)),

		// --- usecase ---
		// todo
		todousecase.NewListUseCase,
		todousecase.NewCreateUseCase,
		todousecase.NewGetUseCase,
		todousecase.NewUpdateUseCase,
		todousecase.NewDeleteUseCase,

		// --- presentation ---
		http.NewRouter,
		handler.NewHandler,
		// health
		health.NewHandler,
		health.NewHealthCheckHandler,
		// todo
		todopresen.NewHandler,
		todopresen.NewListHandler,
		todopresen.NewCreateHandler,
		todopresen.NewGetHandler,
		todopresen.NewUpdateHandler,
		todopresen.NewDeleteHandler,
	)

	// ここはwire generateによって置き換えられる
	return nil, nil
}

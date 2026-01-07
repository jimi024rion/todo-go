//go:generate go tool kessoku $GOFILE

package di

import (
	"github.com/gin-gonic/gin"
	todo_1 "github.com/jimi024rion/todo-go/backend/internal/domain/todo/repository"
	"github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb"
	"github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb/repository/todo"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/handler"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/health"
	todopresen "github.com/jimi024rion/todo-go/backend/internal/presentation/http/todo"
	todousecase "github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
	"github.com/mazrean/kessoku"
)

var _ = kessoku.Inject[*gin.Engine](
	"InitializeServer",
	kessoku.Async(kessoku.Async(kessoku.Provide(rdb.NewTxManager))),
	kessoku.Bind[todo_1.TodoRepository](kessoku.Async(kessoku.Provide(todo.NewRepository))),
	kessoku.Async(kessoku.Provide(todousecase.NewListUseCase)),
	kessoku.Async(kessoku.Provide(todousecase.NewCreateUseCase)),
	kessoku.Async(kessoku.Provide(todousecase.NewGetUseCase)),
	kessoku.Async(kessoku.Provide(todousecase.NewUpdateUseCase)),
	kessoku.Async(kessoku.Provide(todousecase.NewDeleteUseCase)),
	kessoku.Async(kessoku.Provide(http.NewRouter)),
	kessoku.Async(kessoku.Provide(handler.NewHandler)),
	kessoku.Provide(health.NewHandler),
	kessoku.Provide(health.NewHealthCheckHandler),
	kessoku.Async(kessoku.Provide(todopresen.NewHandler)),
	kessoku.Async(kessoku.Provide(todopresen.NewListHandler)),
	kessoku.Async(kessoku.Provide(todopresen.NewCreateHandler)),
	kessoku.Async(kessoku.Provide(todopresen.NewGetHandler)),
	kessoku.Async(kessoku.Provide(todopresen.NewUpdateHandler)),
	kessoku.Async(kessoku.Provide(todopresen.NewDeleteHandler)),
)

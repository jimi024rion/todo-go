//go:generate go tool kessoku $GOFILE

package di

import (
	"github.com/gin-gonic/gin"
	// repository
	"github.com/jimi024rion/todo-go/backend/internal/domain/clock"
	todorepo "github.com/jimi024rion/todo-go/backend/internal/domain/todo/repository"
	clockimpl "github.com/jimi024rion/todo-go/backend/internal/infrastructure/clock"
	"github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb"
	"github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb/repository/todo"

	// presentation
	"github.com/jimi024rion/todo-go/backend/internal/presentation/handler"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http"
	healthpresen "github.com/jimi024rion/todo-go/backend/internal/presentation/http/health"
	todopresen "github.com/jimi024rion/todo-go/backend/internal/presentation/http/todo"
	userpresen "github.com/jimi024rion/todo-go/backend/internal/presentation/http/user"

	// usecase
	todouc "github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
	useruc "github.com/jimi024rion/todo-go/backend/internal/usecase/user"
	"github.com/mazrean/kessoku"
)

var _ = kessoku.Inject[*gin.Engine](
	"InitializeServer",

	// -- Infrastructure Layer --
	// clock
	kessoku.Bind[clock.Clock](kessoku.Provide(clockimpl.NewRealClock)),
	// todo
	kessoku.Bind[todorepo.TodoRepository](kessoku.Async(kessoku.Provide(todo.NewRepository))),

	kessoku.Async(kessoku.Async(kessoku.Provide(rdb.NewTxManager))),

	// -- Usecase Layer --
	// todo
	kessoku.Provide(todouc.NewListUseCase),
	kessoku.Provide(todouc.NewCreateUseCase),
	kessoku.Provide(todouc.NewGetUseCase),
	kessoku.Provide(todouc.NewUpdateUseCase),
	kessoku.Provide(todouc.NewDeleteUseCase),
	// user
	kessoku.Provide(useruc.NewCreateUserUsecase),

	// -- Presentation Layer --
	// health
	kessoku.Provide(healthpresen.NewHandler),
	kessoku.Provide(healthpresen.NewHealthCheckHandler),
	// user
	kessoku.Provide(userpresen.NewCreateUserHandler),
	kessoku.Provide(userpresen.NewHandler),
	// todo
	kessoku.Provide(todopresen.NewListHandler),
	kessoku.Provide(todopresen.NewCreateHandler),
	kessoku.Provide(todopresen.NewGetHandler),
	kessoku.Provide(todopresen.NewUpdateHandler),
	kessoku.Provide(todopresen.NewDeleteHandler),
	kessoku.Provide(todopresen.NewHandler),

	kessoku.Provide(handler.NewHandler),
	kessoku.Provide(http.NewRouter),
)

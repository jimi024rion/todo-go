//go:generate go tool kessoku $GOFILE

package di

import (
	"github.com/gin-gonic/gin"
	"github.com/mazrean/kessoku"

	// repository
	"github.com/jimi024rion/todo-go/backend/internal/config/clock"
	apikeycache "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/cache"
	apikeyrepo "github.com/jimi024rion/todo-go/backend/internal/domain/apikey/repository"
	tagrepo "github.com/jimi024rion/todo-go/backend/internal/domain/tag/repository"
	todorepo "github.com/jimi024rion/todo-go/backend/internal/domain/todo/repository"
	userrepo "github.com/jimi024rion/todo-go/backend/internal/domain/user/repository"
	localCache "github.com/jimi024rion/todo-go/backend/internal/infrastructure/cache/local"
	clockimpl "github.com/jimi024rion/todo-go/backend/internal/infrastructure/clock"
	"github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb"
	apikeyinfra "github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb/repository/apikey"
	taginfra "github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb/repository/tag"
	"github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb/repository/todo"
	userinfra "github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb/repository/user"

	// presentation
	"github.com/jimi024rion/todo-go/backend/internal/presentation/handler"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http"
	apikeypresen "github.com/jimi024rion/todo-go/backend/internal/presentation/http/apikey"
	healthpresen "github.com/jimi024rion/todo-go/backend/internal/presentation/http/health"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/middleware"
	tagpresen "github.com/jimi024rion/todo-go/backend/internal/presentation/http/tag"
	todopresen "github.com/jimi024rion/todo-go/backend/internal/presentation/http/todo"
	userpresen "github.com/jimi024rion/todo-go/backend/internal/presentation/http/user"

	// usecase
	apikeyuc "github.com/jimi024rion/todo-go/backend/internal/usecase/apikey"
	taguc "github.com/jimi024rion/todo-go/backend/internal/usecase/tag"
	todouc "github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
	useruc "github.com/jimi024rion/todo-go/backend/internal/usecase/user"
)

var _ = kessoku.Inject[*gin.Engine](
	"InitializeServer",

	// -- Infrastructure Layer --
	// clock
	kessoku.Bind[clock.Clock](kessoku.Provide(clockimpl.NewRealClock)),
	// todo
	kessoku.Bind[todorepo.TodoRepository](kessoku.Async(kessoku.Provide(todo.NewRepository))),
	// user
	kessoku.Bind[userrepo.UserRepository](kessoku.Async(kessoku.Provide(userinfra.NewRepository))),
	// apikey
	kessoku.Bind[apikeyrepo.APIKeyRepository](kessoku.Async(kessoku.Provide(apikeyinfra.NewRepository))),
	kessoku.Bind[apikeycache.APIKeyCache](kessoku.Provide(localCache.NewLocalAPIKeyCache)),
	// tag
	kessoku.Bind[tagrepo.TagRepository](kessoku.Async(kessoku.Provide(taginfra.NewRepository))),

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
	// apikey
	kessoku.Provide(apikeyuc.NewCreateUseCase),
	kessoku.Provide(apikeyuc.NewDeleteUseCase),
	// tag
	kessoku.Provide(taguc.NewCreateUseCase),
	kessoku.Provide(taguc.NewListUseCase),
	kessoku.Provide(taguc.NewDeleteUseCase),
	kessoku.Provide(taguc.NewSetTodoTagsUseCase),

	// -- Presentation Layer --
	// health
	kessoku.Provide(healthpresen.NewHandler),
	kessoku.Provide(healthpresen.NewHealthCheckHandler),
	// user
	kessoku.Provide(userpresen.NewCreateUserHandler),
	kessoku.Provide(userpresen.NewHandler),
	// apikey
	kessoku.Provide(apikeypresen.NewCreateHandler),
	kessoku.Provide(apikeypresen.NewDeleteHandler),
	kessoku.Provide(apikeypresen.NewHandler),
	// todo
	kessoku.Provide(todopresen.NewListHandler),
	kessoku.Provide(todopresen.NewCreateHandler),
	kessoku.Provide(todopresen.NewGetHandler),
	kessoku.Provide(todopresen.NewUpdateHandler),
	kessoku.Provide(todopresen.NewDeleteHandler),
	kessoku.Provide(todopresen.NewHandler),
	// tag
	kessoku.Provide(tagpresen.NewListHandler),
	kessoku.Provide(tagpresen.NewCreateHandler),
	kessoku.Provide(tagpresen.NewDeleteHandler),
	kessoku.Provide(tagpresen.NewSetTodoTagsHandler),
	kessoku.Provide(tagpresen.NewHandler),

	kessoku.Provide(handler.NewHandler),
	kessoku.Provide(middleware.NewMiddleware),
	kessoku.Provide(http.NewRouter),
)

package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/jimi024rion/todo-go/backend/internal/config/env"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/handler"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/middleware"
)

// NewRouter creates and configures a new Gin router.
func NewRouter(h *handler.Handler, mw *middleware.Middleware) *gin.Engine {
	r := gin.New()

	// Global Middleware
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))
	r.Use(middleware.Logger())
	r.Use(middleware.Trace())
	r.Use(middleware.Timezone())

	// Routes
	r.GET("/health", h.Health.Check.Handle)

	// API v1 group
	v1 := r.Group("/v1")
	{
		// 認証不要
		v1.POST("/api-keys", h.APIKey.Create.Handle)
		v1.POST("/users", h.User.Create.Handle)

		// 認証必須グループ
		auth := v1.Group("")
		auth.Use(mw.FirebaseAuth)
		{
			auth.DELETE("/api-keys/:id", h.APIKey.Delete.Handle)

			todos := auth.Group("/todos")
			{
				todos.GET("", h.Todo.List.Handle)
				todos.POST("", h.Todo.Create.Handle)
				todos.GET("/:id", h.Todo.Get.Handle)
				todos.PUT("/:id", h.Todo.Update.Handle)
				todos.DELETE("/:id", h.Todo.Delete.Handle)
				todos.PUT("/:id/tags", h.Tag.SetTodoTags.Handle)
			}

			tags := auth.Group("/tags")
			{
				tags.GET("", h.Tag.List.Handle)
				tags.POST("", h.Tag.Create.Handle)
				tags.DELETE("/:id", h.Tag.Delete.Handle)
			}
		}
	}

	// Swagger UI (production以外でのみ公開)
	if env.Cfg.AppEnv != "production" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return r
}

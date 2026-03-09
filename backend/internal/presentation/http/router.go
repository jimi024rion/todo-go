package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/handler"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/middleware"
)

// NewRouter creates and configures a new Gin router.
func NewRouter(h *handler.Handler) *gin.Engine {
	r := gin.New()

	// Global Middleware
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.Trace())
	r.Use(middleware.Timezone())

	// Routes
	r.GET("/health", h.Health.Check.Handle)

	// API v1 group
	v1 := r.Group("/v1")
	{

		todos := v1.Group("/todos")
		{
			todos.GET("", h.Todo.List.Handle)
			todos.POST("", h.Todo.Create.Handle)
			todos.GET("/:id", h.Todo.Get.Handle)
			todos.PUT("/:id", h.Todo.Update.Handle)
			todos.DELETE("/:id", h.Todo.Delete.Handle)
		}
	}

	return r
}

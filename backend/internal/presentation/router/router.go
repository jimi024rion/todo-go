package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/internal/presentation/handler"
)

// NewRouter creates and configures a new gin router.
func NewRouter(todoHandler *handler.TodoHandler) *gin.Engine {
	// gin.New() provides more control than gin.Default() if you want to customize middleware.
	// For now, Default() is fine as it includes logger and recovery middleware.
	router := gin.Default()

	// A simple health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Group API endpoints under /api/v1
	apiV1 := router.Group("/api/v1")
	{
		todos := apiV1.Group("/todos")
		{
			todos.POST("", todoHandler.Create)
			todos.GET("", todoHandler.GetAll)
			todos.GET("/:id", todoHandler.GetByID)
			todos.PUT("/:id", todoHandler.Update)
			todos.DELETE("/:id", todoHandler.Delete)
		}
	}

	return router
}

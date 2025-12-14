package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/middleware"
)

// NewRouter creates and configures a new Gin router.
func NewRouter() *gin.Engine {
	r := gin.New()

	// Global Middleware
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	// Routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API v1 group
	// v1 := r.Group("/v1")
	// {
	// 	// Register more routes here
	// }

	return r
}

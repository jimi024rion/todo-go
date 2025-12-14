package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/middleware"
)

// NewRouter creates and configures a new Gin router.
func NewRouter(h *Handler) *gin.Engine {
	r := gin.New()

	// Global Middleware
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	// Routes
	r.GET("/health", h.Health.Check)

	// API v1 group
	// v1 := r.Group("/v1")
	// {
	// 	// Register more routes here
	// }

	return r
}

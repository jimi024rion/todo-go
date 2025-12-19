package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/config/logger"
)

type HealthCheckHandler struct{}

// NewHealthCheckHandler はDIのためのコンストラクタです。
func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

// Handle は GET /health のリクエストを処理します。
func (h *HealthCheckHandler) Handle(c *gin.Context) {
	l := logger.NewLogger(c)
	l.InfoLog("Health check endpoint accessed")

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

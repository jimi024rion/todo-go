package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthCheckHandler struct{}

// NewHealthCheckHandler はDIのためのコンストラクタです。
func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

// Handle は GET /health のリクエストを処理します。
func (h *HealthCheckHandler) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

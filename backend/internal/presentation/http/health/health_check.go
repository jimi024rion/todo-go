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

// HealthCheckResponse はヘルスチェックのレスポンス。
type HealthCheckResponse struct {
	Status string `json:"status" example:"ok"`
} // @name HealthCheckResponse

// Handle は GET /health のリクエストを処理します。
//
// @Summary     ヘルスチェック
// @Description サーバーの稼働状況を確認します
// @Tags        health
// @Produce     json
// @Success     200 {object} HealthCheckResponse
// @Router      /health [get]
func (h *HealthCheckHandler) Handle(c *gin.Context) {
	l := logger.NewLogger(c)
	l.InfoLog("Health check endpoint accessed")

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

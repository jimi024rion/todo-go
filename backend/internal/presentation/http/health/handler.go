package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler は health パッケージのハンドラ構造体です。
// DB接続確認などが必要な場合は、ここに依存を追加します。
type Handler struct {
	// e.g. db *sql.DB
}

// NewHandler はDIのためのコンストラクタです。
func NewHandler() *Handler {
	return &Handler{}
}

// Check は GET /health のリクエストを処理します。
func (h *Handler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

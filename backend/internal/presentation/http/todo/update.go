package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	todousecase "github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
)

// UpdateHandler is a handler for updating a todo.
type UpdateHandler struct {
	usecase *todousecase.UpdateUseCase
}

// NewUpdateHandler creates a new UpdateHandler.
func NewUpdateHandler(u *todousecase.UpdateUseCase) *UpdateHandler {
	return &UpdateHandler{usecase: u}
}

// updateRequestBody is the request body for updating a todo.
type updateRequestBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// Handle handles the request to update an existing todo.
func (h *UpdateHandler) Handle(c *gin.Context) {
	// URLパラメータからIDを取得
	id := c.Param("id")

	// リクエストボディをバインド
	var req updateRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	// ユースケースの入力DTOを作成
	input := &todousecase.UpdateInput{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}

	// ユースケースを実行
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		if errs.IsBadRequest(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errs.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// 成功レスポンス
	c.JSON(http.StatusOK, output)
}

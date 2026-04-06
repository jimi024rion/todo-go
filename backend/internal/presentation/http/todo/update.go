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
	Title       string `json:"title"       example:"買い物リストを作る"`
	Description string `json:"description" example:"牛乳、卵、パンを買う"`
	Status      string `json:"status"      example:"in_progress" enums:"pending,in_progress,completed"`
}

// Handle handles the request to update an existing todo.
//
// @Summary     タスク更新
// @Description 指定されたIDのタスクを更新します
// @Tags        todos
// @Accept      json
// @Produce     json
// @Param       id      path     string            true "タスクID" example("01234567-89ab-cdef-0123-456789abcdef")
// @Param       request body     updateRequestBody true "タスク更新リクエスト"
// @Success     200     {object} TodoResponse
// @Failure     400     {object} response.ErrorResponse
// @Failure     404     {object} response.ErrorResponse
// @Failure     500     {object} response.ErrorResponse
// @Security    ApiKeyAuth
// @Router      /v1/todos/{id} [put]
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

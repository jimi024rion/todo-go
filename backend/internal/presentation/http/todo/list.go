package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	todousecase "github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
)

// ListHandler is a handler for listing todos.
type ListHandler struct {
	usecase *todousecase.ListUseCase
}

// NewListHandler creates a new ListHandler.
func NewListHandler(u *todousecase.ListUseCase) *ListHandler {
	return &ListHandler{usecase: u}
}

// Handle handles the request to list all todos.
func (h *ListHandler) Handle(c *gin.Context) {
	// ユースケースの入力DTOを作成（現在は空）
	input := &todousecase.ListInput{}

	// ユースケースを実行
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// 成功レスポンス (Todosのスライスを直接返す)
	c.JSON(http.StatusOK, output.Todos)
}

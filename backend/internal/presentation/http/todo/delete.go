package todo

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	todousecase "github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
)

// DeleteHandler is a handler for deleting a todo.
type DeleteHandler struct {
	usecase *todousecase.DeleteUseCase
}

// NewDeleteHandler creates a new DeleteHandler.
func NewDeleteHandler(u *todousecase.DeleteUseCase) *DeleteHandler {
	return &DeleteHandler{usecase: u}
}

// Handle handles the request to delete a todo.
func (h *DeleteHandler) Handle(c *gin.Context) {
	id := c.Param("id")

	input := &todousecase.DeleteInput{
		ID: id,
	}

	_, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		// IDフォーマット不正(UUID形式でない)の場合
		if strings.Contains(err.Error(), "failed to parse todo id") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 本来はリポジトリ層で定義したErrNotFoundで404を返すべきだが、
		// 現状はユースケースから返るエラーの区別が難しいため、
		// 暫定的にInternalServerErrorとしておくのが安全。
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

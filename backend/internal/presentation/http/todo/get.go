package todo

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	todousecase "github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
)

// GetHandler is a handler for getting a todo.
type GetHandler struct {
	usecase *todousecase.GetUseCase
}

// NewGetHandler creates a new GetHandler.
func NewGetHandler(u *todousecase.GetUseCase) *GetHandler {
	return &GetHandler{usecase: u}
}

// Handle handles the request to get a single todo by ID.
func (h *GetHandler) Handle(c *gin.Context) {
	id := c.Param("id")

	input := &todousecase.GetInput{
		ID: id,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		// ユースケースから返るエラーをハンドリング
		// IDフォーマット不正(UUID形式でない)の場合
		if strings.Contains(err.Error(), "failed to parse todo id") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 本来はリポジトリ層で定義したErrNotFoundで404を返すべきだが、
		// 現状はユースケースから返るエラーの区別が難しいため、
		// 暫定的にInternalServerErrorとしておくのが安全。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/config/logger"
	"go.opentelemetry.io/otel"

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
	ctx, span := otel.Tracer("handler").Start(c, "Span1")
	defer span.End()

	l := logger.NewLogger(ctx)
	l.InfoLog("span1-1")
	l.InfoLog("span1-2")
	// ユースケースの入力DTOを作成（現在は空）
	input := &todousecase.ListInput{}

	// ユースケースを実行
	output, err := h.usecase.Execute(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// 成功レスポンス (Todosのスライスを直接返す)
	c.JSON(http.StatusOK, output.Todos)
}

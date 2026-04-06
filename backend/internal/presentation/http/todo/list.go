package todo

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/config/logger"
	"go.opentelemetry.io/otel"

	todousecase "github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
)

// TodoResponse はタスクのAPIレスポンス。
// get, update, list で共通して使用します。
type TodoResponse struct {
	ID          string    `json:"ID"          example:"01234567-89ab-cdef-0123-456789abcdef"`
	Title       string    `json:"Title"       example:"買い物リストを作る"`
	Description string    `json:"Description" example:"牛乳、卵、パンを買う"`
	Status      string    `json:"Status"      example:"pending"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}

// ListHandler is a handler for listing todos.
type ListHandler struct {
	usecase *todousecase.ListUseCase
}

// NewListHandler creates a new ListHandler.
func NewListHandler(u *todousecase.ListUseCase) *ListHandler {
	return &ListHandler{usecase: u}
}

// Handle handles the request to list all todos.
//
// @Summary     タスク一覧取得
// @Description すべてのタスクを一覧で取得します
// @Tags        todos
// @Produce     json
// @Success     200 {array}  TodoResponse
// @Failure     500 {object} response.ErrorResponse
// @Security    ApiKeyAuth
// @Router      /v1/todos [get]
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

package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"

	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	"github.com/jimi024rion/todo-go/backend/internal/config/logger"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/response"
	todousecase "github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
)

// todoListResponse は GET /v1/todos の正常系レスポンスです。
type todoListResponse struct {
	ResultHeader response.ResultHeader `json:"resultHeader"`
	ResultBody   []todoItem            `json:"resultBody"`
} // @name TodoListResponse

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
// @Success     200 {object} todoListResponse
// @Failure     500 {object} response.ErrorNullResponse
// @Security    ApiKeyAuth
// @Router      /v1/todos [get]
func (h *ListHandler) Handle(c *gin.Context) {
	ctx, span := otel.Tracer("handler").Start(c, "Span1")
	defer span.End()

	l := logger.NewLogger(ctx)
	l.InfoLog("span1-1")
	l.InfoLog("span1-2")

	output, err := h.usecase.Execute(ctx, &todousecase.ListInput{})
	if err != nil {
		rc := errs.ResultCodeFrom(err)
		c.JSON(rc.HTTPStatus(), response.FailNull(rc))
		return
	}

	items := make([]todoItem, len(output.Todos))
	for i, t := range output.Todos {
		items[i] = todoItem{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Status:      t.Status,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, todoListResponse{
		ResultHeader: response.OKHeader(),
		ResultBody:   items,
	})
}

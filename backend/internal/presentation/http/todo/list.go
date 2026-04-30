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

type ListHandler struct {
	usecase *todousecase.ListUseCase
}

func NewListHandler(u *todousecase.ListUseCase) *ListHandler {
	return &ListHandler{usecase: u}
}

// Handle handles the request to list all todos.
//
// @Summary     タスク一覧取得
// @Description すべてのタスクを一覧で取得します
// @Tags        todos
// @Produce     json
// @Success     200 {array}  todoItem
// @Failure     500 {object} response.ErrorResponse
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
		c.JSON(rc.HTTPStatus(), response.Fail(rc.Message()))
		return
	}

	items := make([]todoItem, len(output.Todos))
	for i, t := range output.Todos {
		tags := make([]todoTag, len(t.Tags))
		for j, tag := range t.Tags {
			tags[j] = todoTag{ID: tag.ID, Name: tag.Name}
		}
		items[i] = todoItem{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Status:      t.Status,
			DueDate:     t.DueDate,
			SortOrder:   t.SortOrder,
			Tags:        tags,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, items)
}

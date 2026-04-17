package todo

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/response"
	todousecase "github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
)

// todoItem は Todo エンドポイントの共通レスポンスボディです。
// list / get / update で参照します。
type todoItem struct {
	ID          string    `json:"ID"          example:"01234567-89ab-cdef-0123-456789abcdef"`
	Title       string    `json:"Title"       example:"買い物リストを作る"`
	Description string    `json:"Description" example:"牛乳、卵、パンを買う"`
	Status      string    `json:"Status"      example:"pending"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
} // @name TodoItem

// todoGetResponse は GET /v1/todos/:id および PUT /v1/todos/:id の正常系レスポンスです。
type todoGetResponse struct {
	ResultHeader response.ResultHeader `json:"resultHeader"`
	ResultBody   todoItem              `json:"resultBody"`
} // @name TodoGetResponse

// GetHandler is a handler for getting a todo.
type GetHandler struct {
	usecase *todousecase.GetUseCase
}

// NewGetHandler creates a new GetHandler.
func NewGetHandler(u *todousecase.GetUseCase) *GetHandler {
	return &GetHandler{usecase: u}
}

// Handle handles the request to get a single todo by ID.
//
// @Summary     タスク詳細取得
// @Description 指定されたIDのタスクを取得します
// @Tags        todos
// @Produce     json
// @Param       id  path     string true "タスクID" example("01234567-89ab-cdef-0123-456789abcdef")
// @Success     200 {object} todoGetResponse
// @Failure     400 {object} response.ErrorNullResponse
// @Failure     404 {object} response.ErrorNullResponse
// @Failure     500 {object} response.ErrorNullResponse
// @Security    ApiKeyAuth
// @Router      /v1/todos/{id} [get]
func (h *GetHandler) Handle(c *gin.Context) {
	id := c.Param("id")

	output, err := h.usecase.Execute(c.Request.Context(), &todousecase.GetInput{ID: id})
	if err != nil {
		rc := errs.ResultCodeFrom(err)
		c.JSON(rc.HTTPStatus(), response.FailNull(rc))
		return
	}

	c.JSON(http.StatusOK, todoGetResponse{
		ResultHeader: response.OKHeader(),
		ResultBody: todoItem{
			ID:          output.ID,
			Title:       output.Title,
			Description: output.Description,
			Status:      output.Status,
			CreatedAt:   output.CreatedAt,
			UpdatedAt:   output.UpdatedAt,
		},
	})
}

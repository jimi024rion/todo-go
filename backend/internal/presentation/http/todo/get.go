package todo

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/response"
	todousecase "github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
)

type todoTag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type todoItem struct {
	ID          string     `json:"id"          example:"01234567-89ab-cdef-0123-456789abcdef"`
	Title       string     `json:"title"       example:"買い物リストを作る"`
	Description string     `json:"description" example:"牛乳、卵、パンを買う"`
	Status      string     `json:"status"      example:"pending"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Tags        []todoTag  `json:"tags"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
} // @name TodoItem

type GetHandler struct {
	usecase *todousecase.GetUseCase
}

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
// @Success     200 {object} todoItem
// @Failure     400 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Security    ApiKeyAuth
// @Router      /v1/todos/{id} [get]
func (h *GetHandler) Handle(c *gin.Context) {
	id := c.Param("id")

	output, err := h.usecase.Execute(c.Request.Context(), &todousecase.GetInput{ID: id})
	if err != nil {
		rc := errs.ResultCodeFrom(err)
		c.JSON(rc.HTTPStatus(), response.Fail(rc.Message()))
		return
	}

	c.JSON(http.StatusOK, todoItem{
		ID:          output.ID,
		Title:       output.Title,
		Description: output.Description,
		Status:      output.Status,
		DueDate:     output.DueDate,
		Tags:        []todoTag{},
		CreatedAt:   output.CreatedAt,
		UpdatedAt:   output.UpdatedAt,
	})
}

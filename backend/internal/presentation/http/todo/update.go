package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/response"
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
	Title       string `json:"title"       maxLength:"100"                       example:"買い物リストを作る"`
	Description string `json:"description" maxLength:"1000"                      example:"牛乳、卵、パンを買う"`
	Status      string `json:"status"      enums:"pending,in_progress,completed" example:"in_progress"`
} // @name TodoUpdateRequest

// Handle handles the request to update an existing todo.
//
// @Summary     タスク更新
// @Description 指定されたIDのタスクを更新します
// @Tags        todos
// @Accept      json
// @Produce     json
// @Param       id      path     string            true "タスクID" example("01234567-89ab-cdef-0123-456789abcdef")
// @Param       request body     updateRequestBody true "タスク更新リクエスト"
// @Success     200     {object} todoGetResponse
// @Failure     400     {object} response.ErrorNullResponse
// @Failure     404     {object} response.ErrorNullResponse
// @Failure     500     {object} response.ErrorNullResponse
// @Security    ApiKeyAuth
// @Router      /v1/todos/{id} [put]
func (h *UpdateHandler) Handle(c *gin.Context) {
	id := c.Param("id")

	var req updateRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(errs.BadRequest.HTTPStatus(), response.FailNull(errs.BadRequest))
		return
	}

	output, err := h.usecase.Execute(c.Request.Context(), &todousecase.UpdateInput{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	})
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

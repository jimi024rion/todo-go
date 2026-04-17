package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/response"
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
//
// @Summary     タスク削除
// @Description 指定されたIDのタスクを削除します
// @Tags        todos
// @Produce     json
// @Param       id  path     string true "タスクID" example("01234567-89ab-cdef-0123-456789abcdef")
// @Success     204
// @Failure     400 {object} response.ErrorNullResponse
// @Failure     404 {object} response.ErrorNullResponse
// @Failure     500 {object} response.ErrorNullResponse
// @Security    ApiKeyAuth
// @Router      /v1/todos/{id} [delete]
func (h *DeleteHandler) Handle(c *gin.Context) {
	id := c.Param("id")

	input := &todousecase.DeleteInput{ID: id}

	_, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		rc := errs.ResultCodeFrom(err)
		c.JSON(rc.HTTPStatus(), response.FailNull(rc))
		return
	}

	c.Status(http.StatusNoContent)
}

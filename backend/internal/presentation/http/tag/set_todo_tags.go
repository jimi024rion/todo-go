package tag

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/response"
	taguc "github.com/jimi024rion/todo-go/backend/internal/usecase/tag"
)

type setTodoTagsRequest struct {
	TagIDs []string `json:"tag_ids"`
}

type SetTodoTagsHandler struct {
	usecase *taguc.SetTodoTagsUseCase
}

func NewSetTodoTagsHandler(u *taguc.SetTodoTagsUseCase) *SetTodoTagsHandler {
	return &SetTodoTagsHandler{usecase: u}
}

// @Summary TodoにタグをセットするPUT /v1/todos/:id/tags
// @Tags    todos
// @Accept  json
// @Param   id      path string            true "TodoID"
// @Param   request body setTodoTagsRequest true "タグIDリスト"
// @Success 204
// @Security ApiKeyAuth
// @Security    SkipAuth
// @Router  /v1/todos/{id}/tags [put]
func (h *SetTodoTagsHandler) Handle(c *gin.Context) {
	todoID := c.Param("id")

	var req setTodoTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(errs.BadRequest.HTTPStatus(), response.Fail(errs.BadRequest.Message()))
		return
	}

	if err := h.usecase.Execute(c.Request.Context(), &taguc.SetTodoTagsInput{
		TodoID: todoID,
		TagIDs: req.TagIDs,
	}); err != nil {
		rc := errs.ResultCodeFrom(err)
		c.JSON(rc.HTTPStatus(), response.Fail(rc.Message()))
		return
	}

	c.Status(http.StatusNoContent)
}

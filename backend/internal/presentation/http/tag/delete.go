package tag

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/response"
	taguc "github.com/jimi024rion/todo-go/backend/internal/usecase/tag"
)

type DeleteHandler struct {
	usecase *taguc.DeleteUseCase
}

func NewDeleteHandler(u *taguc.DeleteUseCase) *DeleteHandler {
	return &DeleteHandler{usecase: u}
}

// @Summary タグ削除
// @Tags    tags
// @Param   id path string true "タグID"
// @Success 204
// @Security ApiKeyAuth
// @Router  /v1/tags/{id} [delete]
func (h *DeleteHandler) Handle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(errs.IsUnauthorizedRequest.HTTPStatus(), response.Fail(errs.IsUnauthorizedRequest.Message()))
		return
	}

	id := c.Param("id")
	if err := h.usecase.Execute(c.Request.Context(), &taguc.DeleteInput{
		ID:     id,
		UserID: userID.(string),
	}); err != nil {
		rc := errs.ResultCodeFrom(err)
		c.JSON(rc.HTTPStatus(), response.Fail(rc.Message()))
		return
	}

	c.Status(http.StatusNoContent)
}

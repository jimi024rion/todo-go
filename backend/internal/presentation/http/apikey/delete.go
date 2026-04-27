package apikey

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/response"
	apikeyuc "github.com/jimi024rion/todo-go/backend/internal/usecase/apikey"
)

type DeleteHandler struct {
	u *apikeyuc.DeleteUseCase
}

func NewDeleteHandler(u *apikeyuc.DeleteUseCase) *DeleteHandler {
	return &DeleteHandler{u: u}
}

// Handle は DELETE /v1/api-keys/:id のリクエストを処理します。
//
// @Summary     APIキー削除
// @Description 指定されたIDのAPIキーを削除します
// @Tags        api-keys
// @Produce     json
// @Param       id  path string true "APIキーID" example("01234567-89ab-cdef-0123-456789abcdef")
// @Success     204
// @Failure     400 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Security    ApiKeyAuth
// @Router      /v1/api-keys/{id} [delete]
func (h *DeleteHandler) Handle(c *gin.Context) {
	id := c.Param("id")

	_, err := h.u.Execute(c.Request.Context(), &apikeyuc.DeleteInput{ID: id})
	if err != nil {
		rc := errs.ResultCodeFrom(err)
		c.JSON(rc.HTTPStatus(), response.Fail(rc.Message()))
		return
	}

	c.Status(http.StatusNoContent)
}

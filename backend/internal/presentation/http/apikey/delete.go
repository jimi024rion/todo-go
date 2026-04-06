package apikey

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	apikeyuc "github.com/jimi024rion/todo-go/backend/internal/usecase/apikey"
)

type DeleteHandler struct {
	u *apikeyuc.DeleteUseCase
}

func NewDeleteHandler(u *apikeyuc.DeleteUseCase) *DeleteHandler {
	return &DeleteHandler{u: u}
}

func (h *DeleteHandler) Handle(c *gin.Context) {
	id := c.Param("id")

	_, err := h.u.Execute(c.Request.Context(), &apikeyuc.DeleteInput{ID: id})
	if err != nil {
		if errs.IsBadRequest(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

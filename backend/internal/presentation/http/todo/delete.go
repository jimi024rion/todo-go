package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
)

// DeleteHandler is a handler for deleting a todo.
type DeleteHandler struct {
	u *todo.DeleteUseCase
}

// NewDeleteHandler creates a new DeleteHandler.
func NewDeleteHandler(u *todo.DeleteUseCase) *DeleteHandler {
	return &DeleteHandler{u: u}
}

// Handle handles the request to delete a todo.
func (h *DeleteHandler) Handle(c *gin.Context) {
	id := c.Param("id")
	err := h.u.Execute(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

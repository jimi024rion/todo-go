package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteHandler is a handler for deleting a todo.
type DeleteHandler struct {
}

// NewDeleteHandler creates a new DeleteHandler.
func NewDeleteHandler() *DeleteHandler {
	return &DeleteHandler{}
}

// Handle handles the request to delete a todo.
func (h *DeleteHandler) Handle(c *gin.Context) {
	// id := c.Param("id")
	c.Status(http.StatusNoContent)
}

package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateHandler is a handler for updating a todo.
type UpdateHandler struct {
}

// NewUpdateHandler creates a new UpdateHandler.
func NewUpdateHandler() *UpdateHandler {
	return &UpdateHandler{}
}

// Handle handles the request to update an existing todo.
func (h *UpdateHandler) Handle(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "UpdateTodo - Not Implemented", "id": id})
}

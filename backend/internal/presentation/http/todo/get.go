package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetHandler is a handler for getting a todo.
type GetHandler struct {
}

// NewGetHandler creates a new GetHandler.
func NewGetHandler() *GetHandler {
	return &GetHandler{}
}

// Handle handles the request to get a single todo by ID.
func (h *GetHandler) Handle(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "GetTodo - Not Implemented", "id": id})
}

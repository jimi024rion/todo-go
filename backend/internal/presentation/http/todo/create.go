package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateHandler is a handler for creating a todo.
type CreateHandler struct {
}

// NewCreateHandler creates a new CreateHandler.
func NewCreateHandler() *CreateHandler {
	return &CreateHandler{}
}

// Handle handles the request to create a new todo.
func (h *CreateHandler) Handle(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "CreateTodo - Not Implemented"})
}

package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListHandler is a handler for listing todos.
type ListHandler struct {
	// ここにユースケースなどの依存関係が将来的に入る
}

// NewListHandler creates a new ListHandler.
func NewListHandler() *ListHandler {
	return &ListHandler{}
}

// Handle handles the request to list all todos.
func (h *ListHandler) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ListTodos - Not Implemented"})
}

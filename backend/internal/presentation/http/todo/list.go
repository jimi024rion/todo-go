package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
)

// ListHandler is a handler for listing todos.
type ListHandler struct {
	u *todo.ListUseCase
}

// NewListHandler creates a new ListHandler.
func NewListHandler(u *todo.ListUseCase) *ListHandler {
	return &ListHandler{u: u}
}

// Handle handles the request to list all todos.
func (h *ListHandler) Handle(c *gin.Context) {
	todos, err := h.u.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

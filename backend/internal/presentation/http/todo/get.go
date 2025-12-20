package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
)

// GetHandler is a handler for getting a todo.
type GetHandler struct {
	u *todo.GetUseCase
}

// NewGetHandler creates a new GetHandler.
func NewGetHandler(u *todo.GetUseCase) *GetHandler {
	return &GetHandler{u: u}
}

// Handle handles the request to get a single todo by ID.
func (h *GetHandler) Handle(c *gin.Context) {
	id := c.Param("id")
	todo, err := h.u.Execute(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

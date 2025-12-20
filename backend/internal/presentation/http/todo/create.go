package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
)

// CreateHandler is a handler for creating a todo.
type CreateHandler struct {
	u *todo.CreateUseCase
}

// NewCreateHandler creates a new CreateHandler.
func NewCreateHandler(u *todo.CreateUseCase) *CreateHandler {
	return &CreateHandler{u: u}
}

// Handle handles the request to create a new todo.
func (h *CreateHandler) Handle(c *gin.Context) {
	var request struct {
		Item string `json:"item"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo, err := h.u.Execute(request.Item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, todo)
}

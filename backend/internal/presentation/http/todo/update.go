package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
)

// UpdateHandler is a handler for updating a todo.
type UpdateHandler struct {
	u *todo.UpdateUseCase
}

// NewUpdateHandler creates a new UpdateHandler.
func NewUpdateHandler(u *todo.UpdateUseCase) *UpdateHandler {
	return &UpdateHandler{u: u}
}

// Handle handles the request to update an existing todo.
func (h *UpdateHandler) Handle(c *gin.Context) {
	id := c.Param("id")
	var request struct {
		Item string `json:"item"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo, err := h.u.Execute(id, request.Item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

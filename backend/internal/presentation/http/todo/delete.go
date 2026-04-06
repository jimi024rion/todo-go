package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	todousecase "github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
)

// DeleteHandler is a handler for deleting a todo.
type DeleteHandler struct {
	usecase *todousecase.DeleteUseCase
}

// NewDeleteHandler creates a new DeleteHandler.
func NewDeleteHandler(u *todousecase.DeleteUseCase) *DeleteHandler {
	return &DeleteHandler{usecase: u}
}

// Handle handles the request to delete a todo.
func (h *DeleteHandler) Handle(c *gin.Context) {
	id := c.Param("id")

	input := &todousecase.DeleteInput{
		ID: id,
	}

	_, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		if errs.IsBadRequest(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errs.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

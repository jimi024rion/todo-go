package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	todousecase "github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
)

// GetHandler is a handler for getting a todo.
type GetHandler struct {
	usecase *todousecase.GetUseCase
}

// NewGetHandler creates a new GetHandler.
func NewGetHandler(u *todousecase.GetUseCase) *GetHandler {
	return &GetHandler{usecase: u}
}

// Handle handles the request to get a single todo by ID.
func (h *GetHandler) Handle(c *gin.Context) {
	id := c.Param("id")

	input := &todousecase.GetInput{
		ID: id,
	}

	output, err := h.usecase.Execute(c.Request.Context(), input)
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

	c.JSON(http.StatusOK, output)
}

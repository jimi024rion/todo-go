package apikey

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	apikeyuc "github.com/jimi024rion/todo-go/backend/internal/usecase/apikey"
)

type CreateHandler struct {
	u *apikeyuc.CreateUseCase
}

func NewCreateHandler(u *apikeyuc.CreateUseCase) *CreateHandler {
	return &CreateHandler{u: u}
}

type createRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Name   string `json:"name"    binding:"required"`
}

type createResponse struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *CreateHandler) Handle(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	output, err := h.u.Execute(c.Request.Context(), &apikeyuc.CreateInput{
		UserID: req.UserID,
		Name:   req.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, createResponse{
		ID:        output.ID,
		Key:       output.Key,
		UserID:    output.UserID,
		Name:      output.Name,
		CreatedAt: output.CreatedAt,
	})
}

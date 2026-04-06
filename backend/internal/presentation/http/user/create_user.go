package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	useruc "github.com/jimi024rion/todo-go/backend/internal/usecase/user"
)

type CreateUserHandler struct {
	u *useruc.CreateUserUsecase
}

func NewCreateUserHandler(createUserUsecase *useruc.CreateUserUsecase) *CreateUserHandler {
	return &CreateUserHandler{
		u: createUserUsecase,
	}
}

type createUserRequest struct {
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required"`
}

type createUserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *CreateUserHandler) Handle(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	output, err := h.u.Execute(c.Request.Context(), &useruc.CreateUserInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		if errs.IsBadRequest(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, createUserResponse{
		ID:        output.ID,
		Name:      output.Name,
		Email:     output.Email,
		CreatedAt: output.CreatedAt,
		UpdatedAt: output.UpdatedAt,
	})
}

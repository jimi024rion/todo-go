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
	Name  string `json:"name"  binding:"required" example:"John Doe"`
	Email string `json:"email" binding:"required" example:"john@example.com"`
}

type createUserResponse struct {
	ID        string    `json:"id"         example:"01234567-89ab-cdef-0123-456789abcdef"`
	Name      string    `json:"name"       example:"John Doe"`
	Email     string    `json:"email"      example:"john@example.com"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Handle は POST /v1/users のリクエストを処理します。
//
// @Summary     ユーザー作成
// @Description 新しいユーザーを作成します
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       request body     createUserRequest true "ユーザー作成リクエスト"
// @Success     201     {object} createUserResponse
// @Failure     400     {object} response.ErrorResponse
// @Failure     500     {object} response.ErrorResponse
// @Security    ApiKeyAuth
// @Router      /v1/users [post]
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

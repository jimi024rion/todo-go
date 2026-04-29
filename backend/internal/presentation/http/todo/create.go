package todo

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/response"
	todousecase "github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
)

type CreateHandler struct {
	usecase *todousecase.CreateUseCase
}

func NewCreateHandler(u *todousecase.CreateUseCase) *CreateHandler {
	return &CreateHandler{usecase: u}
}

type createRequestBody struct {
	Title       string     `json:"title"       validate:"required" minLength:"1" maxLength:"100"  example:"買い物リストを作る"`
	Description string     `json:"description"                     maxLength:"1000"               example:"牛乳、卵、パンを買う"`
	DueDate     *time.Time `json:"due_date,omitempty"`
} // @name TodoCreateRequest

type todoCreateBody struct {
	ID string `json:"id" example:"01234567-89ab-cdef-0123-456789abcdef"`
} // @name TodoCreateBody

func (h *CreateHandler) Handle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(errs.IsUnauthorizedRequest.HTTPStatus(), response.Fail(errs.IsUnauthorizedRequest.Message()))
		return
	}

	var req createRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(errs.BadRequest.HTTPStatus(), response.Fail(errs.BadRequest.Message()))
		return
	}

	output, err := h.usecase.Execute(c.Request.Context(), &todousecase.CreateInput{
		UserID:      userID.(string),
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
	})
	if err != nil {
		rc := errs.ResultCodeFrom(err)
		c.JSON(rc.HTTPStatus(), response.Fail(rc.Message()))
		return
	}

	c.JSON(http.StatusCreated, todoCreateBody{ID: output.ID})
}

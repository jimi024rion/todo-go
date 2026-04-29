package tag

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/response"
	taguc "github.com/jimi024rion/todo-go/backend/internal/usecase/tag"
)

type createRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateHandler struct {
	usecase *taguc.CreateUseCase
}

func NewCreateHandler(u *taguc.CreateUseCase) *CreateHandler {
	return &CreateHandler{usecase: u}
}

// @Summary タグ作成
// @Tags    tags
// @Accept  json
// @Produce json
// @Param   request body     createRequest true "タグ作成リクエスト"
// @Success 201     {object} tagItem
// @Security ApiKeyAuth
// @Router  /v1/tags [post]
func (h *CreateHandler) Handle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(errs.IsUnauthorizedRequest.HTTPStatus(), response.Fail(errs.IsUnauthorizedRequest.Message()))
		return
	}

	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(errs.BadRequest.HTTPStatus(), response.Fail(errs.BadRequest.Message()))
		return
	}

	output, err := h.usecase.Execute(c.Request.Context(), &taguc.CreateInput{
		UserID: userID.(string),
		Name:   req.Name,
	})
	if err != nil {
		rc := errs.ResultCodeFrom(err)
		c.JSON(rc.HTTPStatus(), response.Fail(rc.Message()))
		return
	}

	c.JSON(http.StatusCreated, tagItem{
		ID:        output.ID,
		Name:      output.Name,
		CreatedAt: output.CreatedAt,
		UpdatedAt: output.UpdatedAt,
	})
}

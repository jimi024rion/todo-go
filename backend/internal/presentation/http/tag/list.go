package tag

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/response"
	taguc "github.com/jimi024rion/todo-go/backend/internal/usecase/tag"
)

type tagItem struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListHandler struct {
	usecase *taguc.ListUseCase
}

func NewListHandler(u *taguc.ListUseCase) *ListHandler {
	return &ListHandler{usecase: u}
}

// @Summary タグ一覧取得
// @Tags    tags
// @Produce json
// @Success 200 {array}  tagItem
// @Security ApiKeyAuth
// @Router  /v1/tags [get]
func (h *ListHandler) Handle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(errs.IsUnauthorizedRequest.HTTPStatus(), response.Fail(errs.IsUnauthorizedRequest.Message()))
		return
	}

	output, err := h.usecase.Execute(c.Request.Context(), &taguc.ListInput{
		UserID: userID.(string),
	})
	if err != nil {
		rc := errs.ResultCodeFrom(err)
		c.JSON(rc.HTTPStatus(), response.Fail(rc.Message()))
		return
	}

	items := make([]tagItem, len(output.Tags))
	for i, t := range output.Tags {
		items[i] = tagItem{
			ID:        t.ID,
			Name:      t.Name,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		}
	}
	c.JSON(http.StatusOK, items)
}

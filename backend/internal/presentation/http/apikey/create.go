package apikey

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/response"
	apikeyuc "github.com/jimi024rion/todo-go/backend/internal/usecase/apikey"
)

type CreateHandler struct {
	u *apikeyuc.CreateUseCase
}

func NewCreateHandler(u *apikeyuc.CreateUseCase) *CreateHandler {
	return &CreateHandler{u: u}
}

type createRequest struct {
	UserID string `json:"user_id" validate:"required" format:"uuid" example:"01234567-89ab-cdef-0123-456789abcdef"`
	Name   string `json:"name"    validate:"required" minLength:"1" maxLength:"50" example:"My API Key"`
} // @name ApiKeyCreateRequest

type apiKeyItem struct {
	ID        string    `json:"id"         example:"01234567-89ab-cdef-0123-456789abcdef"`
	Key       string    `json:"key"        example:"todo_abc123..."`
	UserID    string    `json:"user_id"    example:"01234567-89ab-cdef-0123-456789abcdef"`
	Name      string    `json:"name"       example:"My API Key"`
	CreatedAt time.Time `json:"created_at"`
} // @name ApiKeyItem

// Handle は POST /v1/api-keys のリクエストを処理します。
//
// @Summary     APIキー作成
// @Description 新しいAPIキーを発行します（認証不要）
// @Tags        api-keys
// @Accept      json
// @Produce     json
// @Param       request body     createRequest true "APIキー作成リクエスト"
// @Success     201     {object} apiKeyItem
// @Failure     400     {object} response.ErrorResponse
// @Failure     500     {object} response.ErrorResponse
// @Router      /v1/api-keys [post]
func (h *CreateHandler) Handle(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(errs.BadRequest.HTTPStatus(), response.Fail(errs.BadRequest.Message()))
		return
	}

	output, err := h.u.Execute(c.Request.Context(), &apikeyuc.CreateInput{
		UserID: req.UserID,
		Name:   req.Name,
	})
	if err != nil {
		rc := errs.ResultCodeFrom(err)
		c.JSON(rc.HTTPStatus(), response.Fail(rc.Message()))
		return
	}

	c.JSON(http.StatusCreated, apiKeyItem{
		ID:        output.ID,
		Key:       output.Key,
		UserID:    output.UserID,
		Name:      output.Name,
		CreatedAt: output.CreatedAt,
	})
}

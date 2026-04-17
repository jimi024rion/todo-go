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

// apiKeyItem は POST /v1/api-keys 正常系の resultBody です。
type apiKeyItem struct {
	ID        string    `json:"id"         example:"01234567-89ab-cdef-0123-456789abcdef"`
	Key       string    `json:"key"        example:"todo_abc123..."`
	UserID    string    `json:"user_id"    example:"01234567-89ab-cdef-0123-456789abcdef"`
	Name      string    `json:"name"       example:"My API Key"`
	CreatedAt time.Time `json:"created_at"`
} // @name ApiKeyItem

type apiKeyCreateResponse = response.Response[apiKeyItem] // @name ApiKeyCreateResponse

type apiKeyItemErr struct {
	ID        string    `json:"id"         example:""`
	Key       string    `json:"key"        example:""`
	UserID    string    `json:"user_id"    example:""`
	Name      string    `json:"name"       example:""`
	CreatedAt time.Time `json:"created_at" example:""`
} // @name ApiKeyItemErr

type apiKeyCreateErrorResponse = response.Response[apiKeyItemErr] // @name ApiKeyCreateErrorResponse

// Handle は POST /v1/api-keys のリクエストを処理します。
//
// @Summary     APIキー作成
// @Description 新しいAPIキーを発行します（認証不要）
// @Tags        api-keys
// @Accept      json
// @Produce     json
// @Param       request body     createRequest true "APIキー作成リクエスト"
// @Success     201     {object} apiKeyCreateResponse
// @Failure     400     {object} response.ErrorNullResponse
// @Failure     500     {object} apiKeyCreateErrorResponse
// @Router      /v1/api-keys [post]
func (h *CreateHandler) Handle(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(errs.BadRequest.HTTPStatus(), response.FailNull(errs.BadRequest))
		return
	}

	output, err := h.u.Execute(c.Request.Context(), &apikeyuc.CreateInput{
		UserID: req.UserID,
		Name:   req.Name,
	})
	if err != nil {
		e := errs.AsErr(err)
		resp := response.NewResponse[apiKeyItem](int(e.ResultCode()), nil)
		c.JSON(e.ResultCode().HTTPStatus(), resp)
		return
	}

	body := apiKeyItem{
		ID:        output.ID,
		Key:       output.Key,
		UserID:    output.UserID,
		Name:      output.Name,
		CreatedAt: output.CreatedAt,
	}
	resp := response.NewResponse(int(errs.ResultOK), &body)
	c.JSON(http.StatusCreated, resp)
}

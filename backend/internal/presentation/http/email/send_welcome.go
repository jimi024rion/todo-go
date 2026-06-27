package email

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/response"
	emailuc "github.com/jimi024rion/todo-go/backend/internal/usecase/email"
)

type SendWelcomeHandler struct {
	usecase *emailuc.SendWelcomeUseCase
}

func NewSendWelcomeHandler(u *emailuc.SendWelcomeUseCase) *SendWelcomeHandler {
	return &SendWelcomeHandler{usecase: u}
}

type sendWelcomeRequest struct {
	To []string `json:"to" validate:"required" example:"user@example.com"`
	CC []string `json:"cc,omitempty"            example:"admin@example.com"`
} // @name EmailSendWelcomeRequest

type sendWelcomeResponse struct {
	Message string `json:"message" example:"Welcome email sent successfully"`
} // @name EmailSendWelcomeResponse

// Handle は POST /v1/emails/welcome のリクエストを処理します。
//
// @Summary     ウェルカムメール送信
// @Description ウェルカムメールを送信します
// @Tags        emails
// @Accept      json
// @Produce     json
// @Param       request body     sendWelcomeRequest true "メール送信リクエスト"
// @Success     200     {object} sendWelcomeResponse
// @Failure     400     {object} response.ErrorResponse
// @Failure     500     {object} response.ErrorResponse
// @Security    ApiKeyAuth
// @Security    SkipAuth
// @Router      /v1/emails/welcome [post]
func (h *SendWelcomeHandler) Handle(c *gin.Context) {
	var req sendWelcomeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(errs.BadRequest.HTTPStatus(), response.Fail(errs.BadRequest.Message()))
		return
	}

	output, err := h.usecase.Execute(c.Request.Context(), &emailuc.SendWelcomeInput{
		To: req.To,
		CC: req.CC,
	})
	if err != nil {
		rc := errs.ResultCodeFrom(err)
		c.JSON(rc.HTTPStatus(), response.Fail(rc.Message()))
		return
	}

	c.JSON(http.StatusOK, sendWelcomeResponse{
		Message: output.Message,
	})
}

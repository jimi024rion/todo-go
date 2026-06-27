package file

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/response"
	fileuc "github.com/jimi024rion/todo-go/backend/internal/usecase/file"
)

type downloadURLResponse struct {
	URL       string `json:"url"`
	ExpiresAt string `json:"expires_at"`
}

type GetDownloadURLHandler struct {
	usecase *fileuc.GetDownloadURLUseCase
}

func NewGetDownloadURLHandler(u *fileuc.GetDownloadURLUseCase) *GetDownloadURLHandler {
	return &GetDownloadURLHandler{usecase: u}
}

// @Summary ファイルダウンロードURL取得
// @Tags    files
// @Produce json
// @Param   path query    string              true "GCSオブジェクトパス"
// @Success 200  {object} downloadURLResponse
// @Security FirebaseAuth
// @Router  /v1/files/download-url [get]
func (h *GetDownloadURLHandler) Handle(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(errs.IsUnauthorizedRequest.HTTPStatus(), response.Fail(errs.IsUnauthorizedRequest.Message()))
		return
	}

	path := c.Query("path")
	if path == "" {
		c.JSON(errs.BadRequest.HTTPStatus(), response.Fail(errs.BadRequest.Message()))
		return
	}

	output, err := h.usecase.Execute(c.Request.Context(), &fileuc.GetDownloadURLInput{
		Path: path,
	})
	if err != nil {
		rc := errs.ResultCodeFrom(err)
		c.JSON(rc.HTTPStatus(), response.Fail(rc.Message()))
		return
	}

	c.JSON(http.StatusOK, downloadURLResponse{
		URL:       output.URL,
		ExpiresAt: output.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z"),
	})
}

package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	"github.com/jimi024rion/todo-go/backend/internal/presentation/http/response"
	todousecase "github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
)

// CreateHandler はTodoを作成するためのHTTPハンドラです。
type CreateHandler struct {
	usecase *todousecase.CreateUseCase
}

// NewCreateHandler は新しいCreateHandlerを生成します。
func NewCreateHandler(u *todousecase.CreateUseCase) *CreateHandler {
	return &CreateHandler{usecase: u}
}

// createRequestBody はTodo作成リクエストのボディを表します。
type createRequestBody struct {
	Title       string `json:"title"       validate:"required" minLength:"1" maxLength:"100"  example:"買い物リストを作る"`
	Description string `json:"description"                     maxLength:"1000"               example:"牛乳、卵、パンを買う"`
} // @name TodoCreateRequest

// todoCreateBody は POST /v1/todos 正常系の resultBody です。
type todoCreateBody struct {
	ID string `json:"id" example:"01234567-89ab-cdef-0123-456789abcdef"`
} // @name TodoCreateBody

// TodoCreateResponse は POST /v1/todos の正常系レスポンスです。
// type TodoCreateResponse struct {
// 	ResultHeader response.ResultHeader `json:"resultHeader"`
// 	ResultBody   todoCreateBody        `json:"resultBody"`
// }

type TodoCreateResponse = response.Response[todoCreateBody] // @name TodoCreateResponse

// Handle はTodo作成リクエストを処理します。
//
// @Summary     タスク作成
// @Description 新しいタスクを作成します
// @Tags        todos
// @Accept      json
// @Produce     json
// @Param       request body     createRequestBody true "タスク作成リクエスト"
// @Success     201     {object} TodoCreateResponse
// @Failure     400     {object} response.ErrorNullResponse
// @Failure     401     {object} response.ErrorNullResponse
// @Failure     500     {object} response.ErrorNullResponse
// @Security    ApiKeyAuth
// @Router      /v1/todos [post]
func (h *CreateHandler) Handle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(errs.IsUnauthorizedRequest.HTTPStatus(), response.FailNull(errs.IsUnauthorizedRequest))
		return
	}

	var req createRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(errs.BadRequest.HTTPStatus(), response.FailNull(errs.BadRequest))
		return
	}

	output, err := h.usecase.Execute(c.Request.Context(), &todousecase.CreateInput{
		UserID:      userID.(string),
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		e := errs.AsErr(err)
		resp := response.NewResponse[todoCreateBody](int(e.ResultCode()), nil)
		c.JSON(e.ResultCode().HTTPStatus(), resp)
		return
	}

	body := todoCreateBody{
		ID: output.ID,
	}
	resp := response.NewResponse[todoCreateBody](int(errs.ResultOK), &body)
	// resp := response.NewResponse[todoCreateBody](http.StatusOK, nil)
	c.JSON(http.StatusCreated, resp)
}

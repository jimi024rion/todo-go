package todo

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	todovo "github.com/jimi024rion/todo-go/backend/internal/domain/todo/model/valueobject"
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

// requestBody はTodo作成リクエストのボディを表します。
type requestBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

const dummyUserID = "4D029D7A-E711-4390-AD84-E6E48B184E1C"

// Handle はTodo作成リクエストを処理します。
func (h *CreateHandler) Handle(c *gin.Context) {
	// リクエストボディをバインド
	var req requestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	// ユースケースの入力DTOを作成
	input := &todousecase.CreateInput{
		UserID:      dummyUserID, // FIXME: 将来的に認証から取得した実際のUserIDに置き換える
		Title:       req.Title,
		Description: req.Description,
	}

	// ユースケースを実行
	output, err := h.usecase.Execute(c.Request.Context(), input)
	if err != nil {
		// エラーの種類に応じてHTTPステータスコードを振り分ける
		if errors.Is(err, todovo.ErrTitleIsEmpty) || errors.Is(err, todovo.ErrTitleTooLong) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// その他の予期せぬエラー
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// 成功レスポンス
	c.JSON(http.StatusCreated, gin.H{
		"id": output.ID,
	})
}

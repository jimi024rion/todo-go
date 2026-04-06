package todo

import (
	"context"
	"time"

	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	todovo "github.com/jimi024rion/todo-go/backend/internal/domain/todo/model/valueobject"
	todorepository "github.com/jimi024rion/todo-go/backend/internal/domain/todo/repository"
)

// GetUseCase は、TodoをIDで取得するためのユースケースです。
type GetUseCase struct {
	todoRepo todorepository.TodoRepository
}

// NewGetUseCase は、GetUseCaseを生成します。
func NewGetUseCase(todoRepo todorepository.TodoRepository) *GetUseCase {
	return &GetUseCase{
		todoRepo: todoRepo,
	}
}

// GetInput は、GetUseCaseの入力です。
type GetInput struct {
	ID string
}

// GetOutput は、GetUseCaseの出力です。
// これはプレゼンテーション層に渡すためのDTO（Data Transfer Object）です。
type GetOutput struct {
	ID          string
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Execute は、ユースケースを実行します。
func (uc *GetUseCase) Execute(ctx context.Context, input *GetInput) (*GetOutput, error) {
	// 文字列のIDから、ドメインモデルのTodoID値オブジェクトを生成します。
	// これにより、IDのフォーマットが不正な場合はこの時点でエラーとして扱えます。
	todoID, err := todovo.TodoIDFromString(input.ID)
	if err != nil {
		return nil, errs.NewErr(errs.BadRequest, err)
	}

	// リポジトリからTodoエンティティを検索します。
	foundTodo, err := uc.todoRepo.FindByID(ctx, todoID)
	if err != nil {
		return nil, err
	}

	// 取得したエンティティの情報を、出力用のDTOに詰め替えて返します。
	// ドメインオブジェクトを直接プレゼンテーション層に返さないのがポイントです。
	return &GetOutput{
		ID:          foundTodo.ID().String(),
		Title:       foundTodo.Title().String(),
		Description: foundTodo.Description(),
		Status:      foundTodo.Status().String(),
		CreatedAt:   foundTodo.CreatedAt(),
		UpdatedAt:   foundTodo.UpdatedAt(),
	}, nil
}

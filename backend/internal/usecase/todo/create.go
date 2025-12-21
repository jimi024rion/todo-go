package todo

import (
	"context"

	todomodel "github.com/jimi024rion/todo-go/backend/internal/domain/model/todo"
	todorepository "github.com/jimi024rion/todo-go/backend/internal/domain/repository/todo"
)

// CreateUseCase は、Todoを作成するためのユースケースです。
type CreateUseCase struct {
	todoRepo todorepository.TodoRepository
}

// NewCreateUseCase は、CreateUseCaseを生成します。
func NewCreateUseCase(todoRepo todorepository.TodoRepository) *CreateUseCase {
	return &CreateUseCase{
		todoRepo: todoRepo,
	}
}

// CreateInput は、CreateUseCaseの入力です。
type CreateInput struct {
	Title       string
	Description string
}

// CreateOutput は、CreateUseCaseの出力です。
type CreateOutput struct {
	ID string
}

// Execute は、ユースケースを実行します。
func (uc *CreateUseCase) Execute(ctx context.Context, input *CreateInput) (*CreateOutput, error) {
	// ドメインモデルのファクトリを呼び出し、エンティティを生成します。
	// この中でビジネスルール（タイトルの長さなど）が検証されます。
	newTodo, err := todomodel.NewTodo(input.Title, input.Description)
	if err != nil {
		// ビジネスルール違反の場合、エラーを返す。
		return nil, err
	}

	// リポジトリにエンティティの永続化を依頼します。
	if err := uc.todoRepo.Save(ctx, newTodo); err != nil {
		return nil, err
	}

	// 結果を出力用の構造体に詰めて返します。
	return &CreateOutput{
		ID: newTodo.ID().String(),
	}, nil
}

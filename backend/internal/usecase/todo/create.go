package todo

import (
	"context"

	todomodel "github.com/jimi024rion/todo-go/backend/internal/domain/model/todo"
	usermodel "github.com/jimi024rion/todo-go/backend/internal/domain/model/user"
	todorepository "github.com/jimi024rion/todo-go/backend/internal/domain/repository/todo"
	"github.com/jimi024rion/todo-go/backend/internal/domain/tx"
)

// CreateUseCase は、Todoを作成するためのユースケースです。
type CreateUseCase struct {
	todoRepo  todorepository.TodoRepository
	txManager tx.TxManager
}

// NewCreateUseCase は、CreateUseCaseを生成します。
func NewCreateUseCase(todoRepo todorepository.TodoRepository, txManager tx.TxManager) *CreateUseCase {
	return &CreateUseCase{
		todoRepo:  todoRepo,
		txManager: txManager,
	}
}

// CreateInput は、CreateUseCaseの入力です。
type CreateInput struct {
	UserID      string
	Title       string
	Description string
}

// CreateOutput は、CreateUseCaseの出力です。
type CreateOutput struct {

	ID string
}

// Execute は、ユースケースを実行します。
func (uc *CreateUseCase) Execute(ctx context.Context, input *CreateInput) (*CreateOutput, error) {
	result, err := uc.txManager.Do(ctx, func(txCtx context.Context) (any, error) {
		userID, err := usermodel.UserIDFromString(input.UserID)
		if err != nil {
			return nil, err
		}

		// ドメインモデルのファクトリを呼び出し、エンティティを生成します。
		// この中でビジネスルール（タイトルの長さなど）が検証されます。
		newTodo, err := todomodel.NewTodo(userID, input.Title, input.Description)
		if err != nil {
			// ビジネスルール違反の場合、エラーを返す。
			return nil, err
		}

		// リポジトリにエンティティの永続化を依頼します。
		if err := uc.todoRepo.Save(txCtx, newTodo); err != nil {
			return nil, err
		}

		// 結果を出力用の構造体に詰めて返します。
		return &CreateOutput{
			ID: newTodo.ID().String(),
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return result.(*CreateOutput), nil
}

package todo

import (
	"context"

	todovo "github.com/jimi024rion/todo-go/backend/internal/domain/todo/model/valueobject"
	todorepository "github.com/jimi024rion/todo-go/backend/internal/domain/todo/repository"
)

// DeleteUseCase は、Todoを削除するためのユースケースです。
type DeleteUseCase struct {
	todoRepo todorepository.TodoRepository
}

// NewDeleteUseCase は、DeleteUseCaseを生成します。
func NewDeleteUseCase(todoRepo todorepository.TodoRepository) *DeleteUseCase {
	return &DeleteUseCase{
		todoRepo: todoRepo,
	}
}

// DeleteInput は、DeleteUseCaseの入力です。
type DeleteInput struct {
	ID string
}

// DeleteOutput は、DeleteUseCaseの出力です。
// 削除処理が成功したことを示すため、フィールドは空です。
type DeleteOutput struct{}

// Execute は、ユースケースを実行します。
func (uc *DeleteUseCase) Execute(ctx context.Context, input *DeleteInput) (*DeleteOutput, error) {
	// 1. IDの検証
	todoID, err := todovo.TodoIDFromString(input.ID)
	if err != nil {
		return nil, err
	}

	// 2. リポジトリに削除を依頼
	if err := uc.todoRepo.Delete(ctx, todoID); err != nil {
		return nil, err
	}

	// 3. 成功を示す空の出力を返す
	return &DeleteOutput{}, nil
}

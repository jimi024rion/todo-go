package todo

import (
	"context"
	"time"

	todomodel "github.com/jimi024rion/todo-go/backend/internal/domain/model/todo"
	todorepository "github.com/jimi024rion/todo-go/backend/internal/domain/repository/todo"
)

// UpdateUseCase は、Todoを更新するためのユースケースです。
type UpdateUseCase struct {
	todoRepo todorepository.TodoRepository
}

// NewUpdateUseCase は、UpdateUseCaseを生成します。
func NewUpdateUseCase(todoRepo todorepository.TodoRepository) *UpdateUseCase {
	return &UpdateUseCase{
		todoRepo: todoRepo,
	}
}

// UpdateInput は、UpdateUseCaseの入力です。
type UpdateInput struct {
	ID          string
	Title       string
	Description string
	Status      string
}

// UpdateOutput は、UpdateUseCaseの出力です。
type UpdateOutput struct {
	ID          string
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Execute は、ユースケースを実行します。
func (uc *UpdateUseCase) Execute(ctx context.Context, input *UpdateInput) (*UpdateOutput, error) {
	// 1. IDの検証
	todoID, err := todomodel.TodoIDFromString(input.ID)
	if err != nil {
		return nil, err
	}

	// 2. リポジトリからエンティティを取得
	targetTodo, err := uc.todoRepo.FindByID(ctx, todoID)
	if err != nil {
		return nil, err
	}

	// 3. エンティティの振る舞いメソッドを呼び出して状態を変更
	//    このメソッドの中で、新しい値に対するビジネスルール検証が行われる
	if err := targetTodo.ChangeTitle(input.Title); err != nil {
		return nil, err
	}
	targetTodo.ChangeDescription(input.Description)

	// Statusの更新
	switch todomodel.Status(input.Status) {
	case todomodel.StatusCompleted:
		targetTodo.MarkAsCompleted()
	case todomodel.StatusPending:
		targetTodo.MarkAsPending()
	case todomodel.StatusInProgress:
		targetTodo.MarkAsInProgress()
	default:
		// 不明なステータスが指定された場合は何もしないか、エラーを返す
		// ここでは何もしないポリシーとする
	}

	// 4. 変更されたエンティティを永続化
	if err := uc.todoRepo.Save(ctx, targetTodo); err != nil {
		return nil, err
	}

	// 5. 更新後の状態を出力DTOに詰めて返す
	return &UpdateOutput{
		ID:          targetTodo.ID().String(),
		Title:       targetTodo.Title().String(),
		Description: targetTodo.Description(),
		Status:      targetTodo.Status().String(),
		CreatedAt:   targetTodo.CreatedAt(),
		UpdatedAt:   targetTodo.UpdatedAt(),
	}, nil
}

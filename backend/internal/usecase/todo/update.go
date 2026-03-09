package todo

import (
	"context"
	"time"

	"github.com/jimi024rion/todo-go/backend/internal/domain/clock"
	todovo "github.com/jimi024rion/todo-go/backend/internal/domain/todo/model/valueobject"
	todorepository "github.com/jimi024rion/todo-go/backend/internal/domain/todo/repository"
)

// UpdateUseCase は、Todoを更新するためのユースケースです。
type UpdateUseCase struct {
	todoRepo todorepository.TodoRepository
	clock    clock.Clock
}

// NewUpdateUseCase は、UpdateUseCaseを生成します。
func NewUpdateUseCase(todoRepo todorepository.TodoRepository, clock clock.Clock) *UpdateUseCase {
	return &UpdateUseCase{
		todoRepo: todoRepo,
		clock:    clock,
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
	todoID, err := todovo.TodoIDFromString(input.ID)
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
	now := uc.clock.Now(ctx)
	if err := targetTodo.ChangeTitle(input.Title, now); err != nil {
		return nil, err
	}
	targetTodo.ChangeDescription(input.Description, now)

	// Statusの更新
	switch todovo.Status(input.Status) {
	case todovo.StatusCompleted:
		targetTodo.MarkAsCompleted(now)
	case todovo.StatusInProgress:
		targetTodo.MarkAsInProgress(now)
	case todovo.StatusPending:
		targetTodo.MarkAsPending(now)
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

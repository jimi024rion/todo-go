package todo

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"

	"github.com/jimi024rion/todo-go/backend/internal/config/clock"
	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	"github.com/jimi024rion/todo-go/backend/internal/domain/todo/model/entity"
	todovo "github.com/jimi024rion/todo-go/backend/internal/domain/todo/model/valueobject"
	todorepository "github.com/jimi024rion/todo-go/backend/internal/domain/todo/repository"
	"github.com/jimi024rion/todo-go/backend/internal/domain/tx"
	uservo "github.com/jimi024rion/todo-go/backend/internal/domain/user/model/valueobject"
)

// CreateUseCase は、Todoを作成するためのユースケースです。
type CreateUseCase struct {
	todoRepo  todorepository.TodoRepository
	txManager tx.TxManager
	clock     clock.Clock
}

// NewCreateUseCase は、CreateUseCaseを生成します。
func NewCreateUseCase(todoRepo todorepository.TodoRepository, txManager tx.TxManager, clock clock.Clock) *CreateUseCase {
	return &CreateUseCase{
		todoRepo:  todoRepo,
		txManager: txManager,
		clock:     clock,
	}
}

// CreateInput は、CreateUseCaseの入力です。
type CreateInput struct {
	UserID      string
	Title       string
	Description string
	DueDate     *time.Time
}

// CreateOutput は、CreateUseCaseの出力です。
type CreateOutput struct {
	ID string
}

// Execute は、ユースケースを実行します。
func (uc *CreateUseCase) Execute(ctx context.Context, input *CreateInput) (*CreateOutput, error) {
	tr := otel.Tracer("todo-usecase")
	ctx, span := tr.Start(ctx, "CreateUseCase.Execute")
	defer span.End()

	result, err := uc.txManager.Do(ctx, func(txCtx context.Context) (any, error) {
		userID, err := uservo.UserIDFromString(input.UserID)
		if err != nil {
			return nil, errs.NewErr(errs.InternalCodeInvalidID, err)
		}

		// ドメインモデルのファクトリを呼び出し、エンティティを生成します。
		// この中でビジネスルール（タイトルの長さなど）が検証されます。
		newTodo, err := entity.NewTodo(userID, input.Title, input.Description, input.DueDate, uc.clock.Now(ctx))
		if err != nil {
			switch {
			case errors.Is(err, todovo.ErrTitleIsEmpty):
				return nil, errs.NewErr(errs.InternalCodeTitleEmpty, err)
			case errors.Is(err, todovo.ErrTitleTooLong):
				return nil, errs.NewErr(errs.InternalCodeTitleTooLong, err)
			default:
				return nil, errs.NewErr(errs.InternalCodeInternal, err)
			}
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

package todo

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"

	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	"github.com/jimi024rion/todo-go/backend/internal/config/logger"
	todorepository "github.com/jimi024rion/todo-go/backend/internal/domain/todo/repository"
)

// ListUseCase は、Todoを一覧取得するためのユースケースです。
type ListUseCase struct {
	todoRepo todorepository.TodoRepository
}

// NewListUseCase は、ListUseCaseを生成します。
func NewListUseCase(todoRepo todorepository.TodoRepository) *ListUseCase {
	return &ListUseCase{
		todoRepo: todoRepo,
	}
}

// ListInput は、ListUseCaseの入力です。
// NOTE: 将来的にページネーションやフィルタリングのパラメータを追加します。
type ListInput struct{}

// ListOutput は、ListUseCaseの出力です。
type ListOutput struct {
	Todos []*TodoOutput
}

// TodoOutput は、一覧取得の各Todo要素を表すDTOです。
type TodoOutput struct {
	ID          string
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Execute は、ユースケースを実行します。
func (uc *ListUseCase) Execute(ctx context.Context, input *ListInput) (*ListOutput, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "Span2")
	defer span.End()

	l := logger.NewLogger(ctx)
	// l.InfoLog("span2-1")
	// l.InfoLog("span2-2")

	l.ErrorLog(errs.NewErr(errs.CodeUnknownErr, errors.New("エラーだよ。スタックトレース出るかなあ")))
	// リポジトリからTodoエンティティのリストを取得します。
	todos, err := uc.todoRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	// 取得したエンティティのリストを、出力用のDTOのリストに詰め替えます。
	todoOutputs := make([]*TodoOutput, len(todos))
	for i, todo := range todos {
		todoOutputs[i] = &TodoOutput{
			ID:          todo.ID().String(),
			Title:       todo.Title().String(),
			Description: todo.Description(),
			Status:      todo.Status().String(),
			CreatedAt:   todo.CreatedAt(),
			UpdatedAt:   todo.UpdatedAt(),
		}
	}

	return &ListOutput{
		Todos: todoOutputs,
	}, nil
}

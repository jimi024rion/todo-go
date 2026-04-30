package todo

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"

	tagentity "github.com/jimi024rion/todo-go/backend/internal/domain/tag/model/entity"
	tagrepo "github.com/jimi024rion/todo-go/backend/internal/domain/tag/repository"
	todorepository "github.com/jimi024rion/todo-go/backend/internal/domain/todo/repository"
)

// ListUseCase は、Todoを一覧取得するためのユースケースです。
type ListUseCase struct {
	todoRepo todorepository.TodoRepository
	tagRepo  tagrepo.TagRepository
}

// NewListUseCase は、ListUseCaseを生成します。
func NewListUseCase(todoRepo todorepository.TodoRepository, tagRepo tagrepo.TagRepository) *ListUseCase {
	return &ListUseCase{todoRepo: todoRepo, tagRepo: tagRepo}
}

// ListInput は、ListUseCaseの入力です。
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
	DueDate     *time.Time
	SortOrder   float64
	Tags        []TagOutput
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TagOutput はタグのDTOです。
type TagOutput struct {
	ID   string
	Name string
}

// Execute は、ユースケースを実行します。
func (uc *ListUseCase) Execute(ctx context.Context, input *ListInput) (*ListOutput, error) {
	ctx, span := otel.Tracer("handler").Start(ctx, "Span2")
	defer span.End()

	todos, err := uc.todoRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	todoIDs := make([]string, len(todos))
	for i, t := range todos {
		todoIDs[i] = t.ID().String()
	}

	tagsByTodoID, err := uc.tagRepo.FindByTodoIDs(ctx, todoIDs)
	if err != nil {
		tagsByTodoID = map[string][]*tagentity.Tag{}
	}

	todoOutputs := make([]*TodoOutput, len(todos))
	for i, todo := range todos {
		todoID := todo.ID().String()
		tagOutputs := toTagOutputs(tagsByTodoID[todoID])
		todoOutputs[i] = &TodoOutput{
			ID:          todoID,
			Title:       todo.Title().String(),
			Description: todo.Description(),
			Status:      todo.Status().String(),
			DueDate:     todo.DueDate(),
			SortOrder:   todo.SortOrder(),
			Tags:        tagOutputs,
			CreatedAt:   todo.CreatedAt(),
			UpdatedAt:   todo.UpdatedAt(),
		}
	}

	return &ListOutput{Todos: todoOutputs}, nil
}

func toTagOutputs(tags []*tagentity.Tag) []TagOutput {
	out := make([]TagOutput, len(tags))
	for i, t := range tags {
		out[i] = TagOutput{ID: t.ID(), Name: t.Name()}
	}
	return out
}

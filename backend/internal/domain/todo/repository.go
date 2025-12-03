package todo

import (
	"context"

	"github.com/jimi024rion/todo-go/internal/domain/todo/entity"
)

// TodoRepository defines the interface for database operations on todos.
type TodoRepository interface {
	FindAll(ctx context.Context) ([]*entity.Todo, error)
	FindByID(ctx context.Context, id int64) (*entity.Todo, error)
	Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
	Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
	Delete(ctx context.Context, id int64) error
}
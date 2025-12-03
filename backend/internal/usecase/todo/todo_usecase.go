package todo

import (
	"context"
	"errors"
	"time"

	"github.com/jimi024rion/todo-go/internal/domain/todo"
	"github.com/jimi024rion/todo-go/internal/domain/todo/entity"
)

// I a custom error type for not found errors
var ErrNotFound = errors.New("not found")

type Usecase interface {
	CreateTodo(ctx context.Context, title, description string) (*entity.Todo, error)
	GetTodoByID(ctx context.Context, id int64) (*entity.Todo, error)
	GetAllTodos(ctx context.Context) ([]*entity.Todo, error)
	UpdateTodo(ctx context.Context, id int64, title *string, description *string, completed *bool) (*entity.Todo, error)
	DeleteTodo(ctx context.Context, id int64) error
}

type usecase struct {
	todoRepo todo.TodoRepository
}

// NewUsecase creates a new todo usecase.
func NewUsecase(todoRepo todo.TodoRepository) Usecase {
	return &usecase{todoRepo: todoRepo}
}

func (uc *usecase) CreateTodo(ctx context.Context, title, description string) (*entity.Todo, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}
	now := time.Now()
	newTodo := &entity.Todo{
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return uc.todoRepo.Create(ctx, newTodo)
}

func (uc *usecase) GetTodoByID(ctx context.Context, id int64) (*entity.Todo, error) {
	todo, err := uc.todoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, ErrNotFound
	}
	return todo, nil
}

func (uc *usecase) GetAllTodos(ctx context.Context) ([]*entity.Todo, error) {
	return uc.todoRepo.FindAll(ctx)
}

func (uc *usecase) UpdateTodo(ctx context.Context, id int64, title *string, description *string, completed *bool) (*entity.Todo, error) {
	targetTodo, err := uc.todoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if targetTodo == nil {
		return nil, ErrNotFound
	}

	if title != nil {
		if *title == "" {
			return nil, errors.New("title cannot be empty")
		}
		targetTodo.Title = *title
	}
	if description != nil {
		targetTodo.Description = *description
	}
	if completed != nil {
		targetTodo.Completed = *completed
	}
	targetTodo.UpdatedAt = time.Now()

	return uc.todoRepo.Update(ctx, targetTodo)
}

func (uc *usecase) DeleteTodo(ctx context.Context, id int64) error {
	err := uc.todoRepo.Delete(ctx, id)
	if err != nil && err.Error() == "sql: no rows in result set" { // A bit brittle, better to use custom errors
		return ErrNotFound
	}
	return err
}

package todo

import (
	"context"
	"testing"

	"github.com/jimi024rion/todo-go/internal/domain/todo/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTodoRepository is a mock type for the TodoRepository.
type MockTodoRepository struct {
	mock.Mock
}

// These are the mock implementations of the repository interface.
func (m *MockTodoRepository) Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	args := m.Called(ctx, todo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Todo), args.Error(1)
}

func (m *MockTodoRepository) FindByID(ctx context.Context, id int64) (*entity.Todo, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Todo), args.Error(1)
}
func (m *MockTodoRepository) FindAll(ctx context.Context) ([]*entity.Todo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Todo), args.Error(1)
}
func (m *MockTodoRepository) Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	args := m.Called(ctx, todo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Todo), args.Error(1)
}
func (m *MockTodoRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateTodo(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	uc := NewUsecase(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		title := "Test Title"
		description := "Test Description"
		
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Todo")).Return(&entity.Todo{ID: 1, Title: title, Description: description}, nil).Once()

		createdTodo, err := uc.CreateTodo(ctx, title, description)

		assert.NoError(t, err)
		assert.NotNil(t, createdTodo)
		assert.Equal(t, title, createdTodo.Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty title", func(t *testing.T) {
		_, err := uc.CreateTodo(ctx, "", "description")
		assert.Error(t, err)
	})
}

func TestGetTodoByID(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	uc := NewUsecase(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedTodo := &entity.Todo{ID: 1, Title: "Test"}
		mockRepo.On("FindByID", ctx, int64(1)).Return(expectedTodo, nil).Once()

		result, err := uc.GetTodoByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedTodo, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, int64(99)).Return(nil, nil).Once()

		_, err := uc.GetTodoByID(ctx, 99)

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

package todo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	todorepo_mocks "github.com/jimi024rion/todo-go/backend/internal/domain/todo/repository/mocks"
	todovo "github.com/jimi024rion/todo-go/backend/internal/domain/todo/model/valueobject"
	todousecase "github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New().String()
	todo := newTodoEntity(userID, "タスク", "説明")

	repoMock := todorepo_mocks.NewMockTodoRepository(t)
	repoMock.EXPECT().FindByID(mock.Anything, todo.ID()).Return(todo, nil)

	uc := todousecase.NewGetUseCase(repoMock)
	output, err := uc.Execute(ctx, &todousecase.GetInput{ID: todo.ID().String()})

	require.NoError(t, err)
	assert.Equal(t, todo.ID().String(), output.ID)
	assert.Equal(t, "タスク", output.Title)
	assert.Equal(t, "説明", output.Description)
}

func TestGetUseCase_Execute_NotFound(t *testing.T) {
	ctx := context.Background()
	notFoundID := uuid.New().String()
	id, _ := todovo.TodoIDFromString(notFoundID)

	repoMock := todorepo_mocks.NewMockTodoRepository(t)
	repoMock.EXPECT().FindByID(mock.Anything, id).
		Return(nil, errs.NewErr(errs.InternalCodeNotFound, assert.AnError))

	uc := todousecase.NewGetUseCase(repoMock)
	_, err := uc.Execute(ctx, &todousecase.GetInput{ID: notFoundID})

	require.Error(t, err)
	assert.True(t, errs.IsNotFound(err))
}

func TestGetUseCase_Execute_InvalidID(t *testing.T) {
	ctx := context.Background()

	repoMock := todorepo_mocks.NewMockTodoRepository(t)
	uc := todousecase.NewGetUseCase(repoMock)

	_, err := uc.Execute(ctx, &todousecase.GetInput{ID: "invalid-uuid"})

	require.Error(t, err)
	assert.True(t, errs.IsBadRequest(err))
}

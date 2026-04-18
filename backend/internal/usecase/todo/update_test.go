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

func TestUpdateUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New().String()
	todo := newTodoEntity(userID, "古いタイトル", "古い説明")

	repoMock := todorepo_mocks.NewMockTodoRepository(t)
	repoMock.EXPECT().FindByID(mock.Anything, todo.ID()).Return(todo, nil)
	repoMock.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

	txMock := newTxMock(t)
	uc := todousecase.NewUpdateUseCase(repoMock, txMock, testClock)

	output, err := uc.Execute(ctx, &todousecase.UpdateInput{
		ID:          todo.ID().String(),
		Title:       "新しいタイトル",
		Description: "新しい説明",
		Status:      "in_progress",
	})

	require.NoError(t, err)
	assert.Equal(t, todo.ID().String(), output.ID)
	assert.Equal(t, "新しいタイトル", output.Title)
	assert.Equal(t, "新しい説明", output.Description)
	assert.Equal(t, "in_progress", output.Status)
}

func TestUpdateUseCase_Execute_EmptyTitle(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New().String()
	todo := newTodoEntity(userID, "タイトル", "")

	repoMock := todorepo_mocks.NewMockTodoRepository(t)
	repoMock.EXPECT().FindByID(mock.Anything, todo.ID()).Return(todo, nil)

	txMock := newTxMock(t)
	uc := todousecase.NewUpdateUseCase(repoMock, txMock, testClock)

	_, err := uc.Execute(ctx, &todousecase.UpdateInput{
		ID:    todo.ID().String(),
		Title: "",
	})

	require.Error(t, err)
	assert.True(t, errs.IsBadRequest(err))
}

func TestUpdateUseCase_Execute_NotFound(t *testing.T) {
	ctx := context.Background()
	notFoundID := uuid.New().String()
	id, _ := todovo.TodoIDFromString(notFoundID)

	repoMock := todorepo_mocks.NewMockTodoRepository(t)
	repoMock.EXPECT().FindByID(mock.Anything, id).
		Return(nil, errs.NewErr(errs.InternalCodeNotFound, assert.AnError))

	txMock := newTxMock(t)
	uc := todousecase.NewUpdateUseCase(repoMock, txMock, testClock)

	_, err := uc.Execute(ctx, &todousecase.UpdateInput{
		ID:    notFoundID,
		Title: "タイトル",
	})

	require.Error(t, err)
	assert.True(t, errs.IsNotFound(err))
}

func TestUpdateUseCase_Execute_InvalidID(t *testing.T) {
	ctx := context.Background()

	repoMock := todorepo_mocks.NewMockTodoRepository(t)
	txMock := newTxMock(t)
	uc := todousecase.NewUpdateUseCase(repoMock, txMock, testClock)

	_, err := uc.Execute(ctx, &todousecase.UpdateInput{
		ID:    "invalid-uuid",
		Title: "タイトル",
	})

	require.Error(t, err)
	assert.True(t, errs.IsBadRequest(err))
}

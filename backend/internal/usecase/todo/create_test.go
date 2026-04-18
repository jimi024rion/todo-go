package todo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	todorepo_mocks "github.com/jimi024rion/todo-go/backend/internal/domain/todo/repository/mocks"
	todousecase "github.com/jimi024rion/todo-go/backend/internal/usecase/todo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New().String()

	repoMock := todorepo_mocks.NewMockTodoRepository(t)
	repoMock.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)

	txMock := newTxMock(t)
	uc := todousecase.NewCreateUseCase(repoMock, txMock, testClock)

	output, err := uc.Execute(ctx, &todousecase.CreateInput{
		UserID:      userID,
		Title:       "テストタスク",
		Description: "テスト説明",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, output.ID)
}

func TestCreateUseCase_Execute_EmptyTitle(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New().String()

	repoMock := todorepo_mocks.NewMockTodoRepository(t)
	txMock := newTxMock(t)
	uc := todousecase.NewCreateUseCase(repoMock, txMock, testClock)

	_, err := uc.Execute(ctx, &todousecase.CreateInput{
		UserID: userID,
		Title:  "",
	})

	require.Error(t, err)
	assert.True(t, errs.IsBadRequest(err))
}

func TestCreateUseCase_Execute_InvalidUserID(t *testing.T) {
	ctx := context.Background()

	repoMock := todorepo_mocks.NewMockTodoRepository(t)
	txMock := newTxMock(t)
	uc := todousecase.NewCreateUseCase(repoMock, txMock, testClock)

	_, err := uc.Execute(ctx, &todousecase.CreateInput{
		UserID: "invalid-uuid",
		Title:  "タイトル",
	})

	require.Error(t, err)
	assert.True(t, errs.IsBadRequest(err))
}

func TestCreateUseCase_Execute_RepoError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New().String()

	repoMock := todorepo_mocks.NewMockTodoRepository(t)
	repoMock.EXPECT().Save(mock.Anything, mock.Anything).
		Return(errs.NewErr(errs.InternalCodeInternal, assert.AnError))

	txMock := newTxMock(t)
	uc := todousecase.NewCreateUseCase(repoMock, txMock, testClock)

	_, err := uc.Execute(ctx, &todousecase.CreateInput{
		UserID: userID,
		Title:  "タイトル",
	})

	require.Error(t, err)
}

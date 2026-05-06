package todo_test

import (
	"context"
	"time"

	"github.com/jimi024rion/todo-go/backend/internal/config/clock"
	"github.com/jimi024rion/todo-go/backend/internal/domain/todo/model/entity"
	txmocks "github.com/jimi024rion/todo-go/backend/internal/domain/tx/mocks"
	uservo "github.com/jimi024rion/todo-go/backend/internal/domain/user/model/valueobject"
	"github.com/stretchr/testify/mock"
)

var (
	testNow   = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	testClock = clock.NewMockClock(testNow)
)

// newTxMock は fn をそのまま実行する TxManager モックを返します。
func newTxMock(t interface {
	mock.TestingT
	Cleanup(func())
},
) *txmocks.MockTxManager {
	txMock := txmocks.NewMockTxManager(t)
	txMock.EXPECT().Do(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) (any, error)) (any, error) {
			return fn(ctx)
		}).
		Maybe()
	return txMock
}

// newTodoEntity はテスト用 Todo エンティティを生成します。
func newTodoEntity(userIDStr, title, description string) *entity.Todo {
	userID, _ := uservo.UserIDFromString(userIDStr)
	todo, _ := entity.NewTodo(userID, title, description, &testNow, testNow)
	return todo
}

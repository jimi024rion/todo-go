package todo_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	"github.com/jimi024rion/todo-go/backend/internal/domain/todo/model/entity"
	todovo "github.com/jimi024rion/todo-go/backend/internal/domain/todo/model/valueobject"
	userentity "github.com/jimi024rion/todo-go/backend/internal/domain/user/model/entity"
	uservo "github.com/jimi024rion/todo-go/backend/internal/domain/user/model/valueobject"
	todorepo "github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb/repository/todo"
	userrepo "github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb/repository/user"
	"github.com/jimi024rion/todo-go/backend/internal/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestUser(t *testing.T, ctx context.Context) *userentity.User {
	t.Helper()
	now := time.Now().Truncate(time.Second)
	user, err := userentity.NewUser(uuid.New().String(), "テストユーザー", uuid.New().String()+"@example.com", now)
	require.NoError(t, err)
	uRepo := userrepo.NewRepository(testBobDB)
	require.NoError(t, uRepo.Save(ctx, user))
	return user
}

func newTestTodo(t *testing.T, ctx context.Context, userID uservo.UserID) *entity.Todo {
	t.Helper()
	now := time.Now().Truncate(time.Second)
	todo, err := entity.NewTodo(userID, "テストタスク", "テスト説明", &now, now)
	require.NoError(t, err)
	repo := todorepo.NewRepository(testBobDB)
	require.NoError(t, repo.Save(ctx, todo))
	return todo
}

func TestTodoRepository_Save_Insert(t *testing.T) {
	ctx := context.Background()
	testhelper.TruncateTables(t, ctx, testSQLDB)

	user := newTestUser(t, ctx)
	repo := todorepo.NewRepository(testBobDB)

	now := time.Now().Truncate(time.Second)
	todo, err := entity.NewTodo(user.ID(), "買い物", "牛乳・卵", &now, now)
	require.NoError(t, err)

	require.NoError(t, repo.Save(ctx, todo))

	found, err := repo.FindByID(ctx, todo.ID())
	require.NoError(t, err)
	assert.Equal(t, todo.ID().String(), found.ID().String())
	assert.Equal(t, "買い物", found.Title().String())
	assert.Equal(t, "牛乳・卵", found.Description())
	assert.Equal(t, todovo.StatusPending, found.Status())
}

func TestTodoRepository_Save_Update(t *testing.T) {
	ctx := context.Background()
	testhelper.TruncateTables(t, ctx, testSQLDB)

	user := newTestUser(t, ctx)
	todo := newTestTodo(t, ctx, user.ID())
	repo := todorepo.NewRepository(testBobDB)

	now := time.Now().Truncate(time.Second)
	require.NoError(t, todo.ChangeTitle("更新後タイトル", now))
	todo.MarkAsCompleted(now)

	require.NoError(t, repo.Save(ctx, todo))

	found, err := repo.FindByID(ctx, todo.ID())
	require.NoError(t, err)
	assert.Equal(t, "更新後タイトル", found.Title().String())
	assert.Equal(t, todovo.StatusCompleted, found.Status())
}

func TestTodoRepository_FindByID_NotFound(t *testing.T) {
	ctx := context.Background()
	testhelper.TruncateTables(t, ctx, testSQLDB)

	repo := todorepo.NewRepository(testBobDB)
	id, err := todovo.TodoIDFromString(uuid.New().String())
	require.NoError(t, err)

	_, err = repo.FindByID(ctx, id)
	require.Error(t, err)
	assert.True(t, errs.IsNotFound(err))
}

func TestTodoRepository_List(t *testing.T) {
	ctx := context.Background()
	testhelper.TruncateTables(t, ctx, testSQLDB)

	user := newTestUser(t, ctx)
	newTestTodo(t, ctx, user.ID())
	newTestTodo(t, ctx, user.ID())

	repo := todorepo.NewRepository(testBobDB)
	todos, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, todos, 2)
}

func TestTodoRepository_List_Empty(t *testing.T) {
	ctx := context.Background()
	testhelper.TruncateTables(t, ctx, testSQLDB)

	repo := todorepo.NewRepository(testBobDB)
	todos, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Empty(t, todos)
}

func TestTodoRepository_Delete(t *testing.T) {
	ctx := context.Background()
	testhelper.TruncateTables(t, ctx, testSQLDB)

	user := newTestUser(t, ctx)
	todo := newTestTodo(t, ctx, user.ID())
	repo := todorepo.NewRepository(testBobDB)

	require.NoError(t, repo.Delete(ctx, todo.ID()))

	_, err := repo.FindByID(ctx, todo.ID())
	require.Error(t, err)
	assert.True(t, errs.IsNotFound(err))
}

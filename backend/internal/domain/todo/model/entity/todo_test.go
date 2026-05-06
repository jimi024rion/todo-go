package entity_test

import (
	"strings"
	"testing"
	"time"

	"github.com/jimi024rion/todo-go/backend/internal/domain/todo/model/entity"
	"github.com/jimi024rion/todo-go/backend/internal/domain/todo/model/valueobject"
	uservo "github.com/jimi024rion/todo-go/backend/internal/domain/user/model/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testUserID = uservo.NewUserID()
	testNow    = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
)

func TestNewTodo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		title       string
		description string
		wantErr     error
	}{
		{
			name:        "正常: タイトルと説明あり",
			title:       "買い物リスト",
			description: "牛乳、卵",
		},
		{
			name:  "正常: 説明なし",
			title: "ミーティング",
		},
		{
			name:    "異常: タイトル空",
			title:   "",
			wantErr: valueobject.ErrTitleIsEmpty,
		},
		{
			name:    "異常: タイトルが100文字超過",
			title:   strings.Repeat("a", 101),
			wantErr: valueobject.ErrTitleTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := entity.NewTodo(testUserID, tt.title, tt.description, &testNow, testNow)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got.ID().String())
			assert.Equal(t, tt.title, got.Title().String())
			assert.Equal(t, tt.description, got.Description())
			assert.Equal(t, valueobject.StatusPending, got.Status())
			assert.Equal(t, testNow, got.CreatedAt())
			assert.Equal(t, testNow, got.UpdatedAt())
		})
	}
}

func TestTodo_MarkAsCompleted(t *testing.T) {
	t.Parallel()

	t.Run("pending→completedに変わりupdatedAtが更新される", func(t *testing.T) {
		t.Parallel()
		todo, err := entity.NewTodo(testUserID, "テスト", "", &testNow, testNow)
		require.NoError(t, err)

		later := testNow.Add(time.Hour)
		todo.MarkAsCompleted(later)

		assert.Equal(t, valueobject.StatusCompleted, todo.Status())
		assert.True(t, todo.IsCompleted())
		assert.Equal(t, later, todo.UpdatedAt())
	})

	t.Run("すでにcompletedなら冪等（updatedAt変わらない）", func(t *testing.T) {
		t.Parallel()
		todo, err := entity.NewTodo(testUserID, "テスト", "", &testNow, testNow)
		require.NoError(t, err)
		todo.MarkAsCompleted(testNow.Add(time.Hour))
		updatedAt := todo.UpdatedAt()

		todo.MarkAsCompleted(testNow.Add(2 * time.Hour))

		assert.Equal(t, updatedAt, todo.UpdatedAt())
	})
}

func TestTodo_MarkAsInProgress(t *testing.T) {
	t.Parallel()

	t.Run("pending→in_progressに変わる", func(t *testing.T) {
		t.Parallel()
		todo, err := entity.NewTodo(testUserID, "テスト", "", &testNow, testNow)
		require.NoError(t, err)

		later := testNow.Add(time.Hour)
		todo.MarkAsInProgress(later)

		assert.Equal(t, valueobject.StatusInProgress, todo.Status())
		assert.Equal(t, later, todo.UpdatedAt())
	})

	t.Run("すでにin_progressなら冪等", func(t *testing.T) {
		t.Parallel()
		todo, err := entity.NewTodo(testUserID, "テスト", "", &testNow, testNow)
		require.NoError(t, err)
		todo.MarkAsInProgress(testNow.Add(time.Hour))
		updatedAt := todo.UpdatedAt()

		todo.MarkAsInProgress(testNow.Add(2 * time.Hour))

		assert.Equal(t, updatedAt, todo.UpdatedAt())
	})
}

func TestTodo_MarkAsPending(t *testing.T) {
	t.Parallel()

	t.Run("completed→pendingに戻る", func(t *testing.T) {
		t.Parallel()
		todo, err := entity.NewTodo(testUserID, "テスト", "", &testNow, testNow)
		require.NoError(t, err)
		todo.MarkAsCompleted(testNow.Add(time.Hour))

		later := testNow.Add(2 * time.Hour)
		todo.MarkAsPending(later)

		assert.Equal(t, valueobject.StatusPending, todo.Status())
		assert.Equal(t, later, todo.UpdatedAt())
	})

	t.Run("すでにpendingなら冪等", func(t *testing.T) {
		t.Parallel()
		todo, err := entity.NewTodo(testUserID, "テスト", "", &testNow, testNow)
		require.NoError(t, err)
		updatedAt := todo.UpdatedAt()

		todo.MarkAsPending(testNow.Add(time.Hour))

		assert.Equal(t, updatedAt, todo.UpdatedAt())
	})
}

func TestTodo_ChangeTitle(t *testing.T) {
	t.Parallel()

	t.Run("正常: タイトルが変わりupdatedAtが更新される", func(t *testing.T) {
		t.Parallel()
		todo, err := entity.NewTodo(testUserID, "古いタイトル", "", &testNow, testNow)
		require.NoError(t, err)

		later := testNow.Add(time.Hour)
		err = todo.ChangeTitle("新しいタイトル", later)

		require.NoError(t, err)
		assert.Equal(t, "新しいタイトル", todo.Title().String())
		assert.Equal(t, later, todo.UpdatedAt())
	})

	t.Run("同じタイトルなら変更なし（updatedAt変わらない）", func(t *testing.T) {
		t.Parallel()
		todo, err := entity.NewTodo(testUserID, "同じタイトル", "", &testNow, testNow)
		require.NoError(t, err)
		updatedAt := todo.UpdatedAt()

		err = todo.ChangeTitle("同じタイトル", testNow.Add(time.Hour))

		require.NoError(t, err)
		assert.Equal(t, updatedAt, todo.UpdatedAt())
	})

	t.Run("異常: 空文字", func(t *testing.T) {
		t.Parallel()
		todo, err := entity.NewTodo(testUserID, "タイトル", "", &testNow, testNow)
		require.NoError(t, err)

		err = todo.ChangeTitle("", testNow.Add(time.Hour))

		require.ErrorIs(t, err, valueobject.ErrTitleIsEmpty)
	})
}

func TestTodo_ChangeDescription(t *testing.T) {
	t.Parallel()

	t.Run("説明が変わりupdatedAtが更新される", func(t *testing.T) {
		t.Parallel()
		todo, err := entity.NewTodo(testUserID, "タイトル", "旧説明", &testNow, testNow)
		require.NoError(t, err)

		later := testNow.Add(time.Hour)
		todo.ChangeDescription("新説明", later)

		assert.Equal(t, "新説明", todo.Description())
		assert.Equal(t, later, todo.UpdatedAt())
	})

	t.Run("同じ説明なら変更なし", func(t *testing.T) {
		t.Parallel()
		todo, err := entity.NewTodo(testUserID, "タイトル", "説明", &testNow, testNow)
		require.NoError(t, err)
		updatedAt := todo.UpdatedAt()

		todo.ChangeDescription("説明", testNow.Add(time.Hour))

		assert.Equal(t, updatedAt, todo.UpdatedAt())
	})
}

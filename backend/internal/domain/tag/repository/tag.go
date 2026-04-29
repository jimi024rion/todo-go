package repository

import (
	"context"

	"github.com/jimi024rion/todo-go/backend/internal/domain/tag/model/entity"
)

// TagRepository はタグの永続化を担うリポジトリインターフェースです。
type TagRepository interface {
	// Create はタグを新規作成します。
	Create(ctx context.Context, tag *entity.Tag) error
	// FindByUserID はユーザーIDに紐づくタグ一覧を返します。
	FindByUserID(ctx context.Context, userID string) ([]*entity.Tag, error)
	// FindByTodoIDs は複数のTodoIDに紐づくタグをまとめて返します。
	// map[todoID][]*Tag の形式で返します。
	FindByTodoIDs(ctx context.Context, todoIDs []string) (map[string][]*entity.Tag, error)
	// Delete はタグを削除します（所有者チェック含む）。
	Delete(ctx context.Context, id, userID string) error
	// SetTodoTags はTodoに紐づくタグを一括更新します（既存は全削除→再挿入）。
	SetTodoTags(ctx context.Context, todoID string, tagIDs []string) error
}

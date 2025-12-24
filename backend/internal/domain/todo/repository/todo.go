package todo

import (
	"context"

	"github.com/jimi024rion/todo-go/backend/internal/domain/todo/model/entity"
	"github.com/jimi024rion/todo-go/backend/internal/domain/todo/model/valueobject"
)

// TodoRepository は Todo エンティティの永続化を抽象化するインターフェースです。
type TodoRepository interface {
	// Save は、Todoエンティティを永続化します。
	// 新規作成か更新かはリポジトリの実装が判断します。
	Save(ctx context.Context, todo *entity.Todo) error

	// FindByID は、指定されたIDを持つTodoエンティティを検索します。
	FindByID(ctx context.Context, id valueobject.TodoID) (*entity.Todo, error)

	// Delete は、指定されたIDを持つTodoエンティティを削除します。
	Delete(ctx context.Context, id valueobject.TodoID) error

	// List は、すべてのTodoエンティティのリストを取得します。
	// TODO: 将来的にページネーションやフィルタリングを考慮する必要があります。
	List(ctx context.Context) ([]*entity.Todo, error)
}

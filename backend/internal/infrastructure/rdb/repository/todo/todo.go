package todo

import (
	"context"
	"fmt"

	"github.com/aarondl/opt/null"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/gofrs/uuid/v5"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/um"

	todomodel "github.com/jimi024rion/todo-go/backend/internal/domain/model/todo"
	"github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb/models"
)

type Repository struct {
	exec bob.Executor
}

// NewRepository は Repository の新しいインスタンスを生成します。
func NewRepository(exec bob.Executor) *Repository {
	return &Repository{exec: exec}
}

// Save は、Todoエンティティを永続化します。
func (r *Repository) Save(ctx context.Context, td *todomodel.Todo) error {
	uid, err := uuid.FromString(td.ID().String())
	if err != nil {
		return fmt.Errorf("invalid todo id format: %w", err)
	}

	exists, err := models.TodoExists(ctx, r.exec, uid)
	if err != nil {
		return fmt.Errorf("failed to check todo existence: %w", err)
	}

	setter := toDBTodoSetter(td)

	if exists {
		// 存在すれば更新
		_, err = models.Todos.Update(
			setter.UpdateMod(),
			um.Where(models.Todos.Columns.ID.EQ(psql.Arg(uid))),
		).Exec(ctx, r.exec)

		if err != nil {
			return fmt.Errorf("failed to update todo: %w", err)
		}
	} else {
		// 存在しなければ挿入
		_, err = models.Todos.Insert(setter).Exec(ctx, r.exec)
		if err != nil {
			return fmt.Errorf("failed to insert todo: %w", err)
		}
	}

	return nil
}

// FindByID は、指定されたIDを持つTodoエンティティを検索します。
func (r *Repository) FindByID(ctx context.Context, id todomodel.TodoID) (*todomodel.Todo, error) {
	uid, err := uuid.FromString(id.String())
	if err != nil {
		return nil, fmt.Errorf("invalid todo id format: %w", err)
	}

	// bobが生成したFindTodoヘルパー関数を利用
	todoModel, err := models.FindTodo(ctx, r.exec, uid)
	if err != nil {
		// TODO: エラーの種類に応じてドメイン層で定義したエラーを返す
		return nil, fmt.Errorf("failed to find todo by id: %w", err)
	}

	return toDomainTodo(todoModel)
}

// Delete は、指定されたIDを持つTodoエンティティを削除します。
func (r *Repository) Delete(ctx context.Context, id todomodel.TodoID) error {
	uid, err := uuid.FromString(id.String())
	if err != nil {
		return fmt.Errorf("invalid todo id format: %w", err)
	}

	_, err = models.Todos.Delete(
		dm.Where(models.Todos.Columns.ID.EQ(psql.Arg(uid))),
	).Exec(ctx, r.exec)

	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	return nil
}

// List は、すべてのTodoエンティティのリストを取得します。
func (r *Repository) List(ctx context.Context) ([]*todomodel.Todo, error) {
	todoModels, err := models.Todos.Query().All(ctx, r.exec)
	if err != nil {
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}

	domainTodos := make([]*todomodel.Todo, len(todoModels))
	for i, m := range todoModels {
		domainTodos[i], err = toDomainTodo(m)
		if err != nil {
			return nil, fmt.Errorf("failed to convert todo model to domain model: %w", err)
		}
	}

	return domainTodos, nil
}

// toDomainTodo は Bob の models.Todo をドメインの todomodel.Todo に変換します。
func toDomainTodo(m *models.Todo) (*todomodel.Todo, error) {
	if m == nil {
		return nil, nil
	}

	description := m.Description.GetOrZero()

	return todomodel.Reconstruct(
		m.ID.String(),
		m.UserID.String(),
		m.Title,
		description,
		m.Status,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

// toDBTodoSetter はドメインの todomodel.Todo を Bob の models.TodoSetter に変換します。
func toDBTodoSetter(td *todomodel.Todo) *models.TodoSetter {
	if td == nil {
		return nil
	}

	var description omitnull.Val[string]
	if td.Description() != "" {
		description = omitnull.From(td.Description())
	} else {
		description = omitnull.FromNull(null.Val[string]{})
	}

	return &models.TodoSetter{
		ID:          omit.From(uuid.FromStringOrNil(td.ID().String())),
		UserID:      omit.From(uuid.FromStringOrNil(td.UserID().String())),
		Title:       omit.From(td.Title().String()),
		Description: description,
		Status:      omit.From(td.Status().String()),
		CreatedAt:   omit.From(td.CreatedAt()),
		UpdatedAt:   omit.From(td.UpdatedAt()),
	}
}

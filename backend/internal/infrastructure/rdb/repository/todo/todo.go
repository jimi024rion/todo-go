package todo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/aarondl/opt/null"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"

	"github.com/jimi024rion/todo-go/backend/internal/config/errs"
	"github.com/jimi024rion/todo-go/backend/internal/domain/todo/model/entity"
	"github.com/jimi024rion/todo-go/backend/internal/domain/todo/model/valueobject"
	"github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb"
	"github.com/jimi024rion/todo-go/backend/internal/infrastructure/rdb/models"
)

type Repository struct {
	db bob.DB
}

// NewRepository は Repository の新しいインスタンスを生成します。
func NewRepository(db bob.DB) *Repository {
	return &Repository{db: db}
}

// getExecutor はコンテキストからトランザクションを取得し、存在すればそれを返します。
// 存在しない場合は、デフォルトのDB接続を返します。
func (r *Repository) getExecutor(ctx context.Context) bob.Executor {
	if exec, ok := rdb.GetExecutor(ctx); ok {
		return exec
	}
	return r.db
}

// Save は、Todoエンティティを永続化します。
func (r *Repository) Save(ctx context.Context, td *entity.Todo) error {
	uid, err := uuid.Parse(td.ID().String())
	if err != nil {
		return fmt.Errorf("invalid todo id format: %w", err)
	}

	exec := r.getExecutor(ctx)
	exists, err := models.TodoExists(ctx, exec, uid)
	if err != nil {
		return fmt.Errorf("failed to check todo existence: %w", err)
	}

	setter := toDBTodoSetter(td)

	if exists {
		// 存在すれば更新
		_, err = models.Todos.Update(
			setter.UpdateMod(),
			um.Where(models.Todos.Columns.ID.EQ(psql.Arg(uid))),
		).Exec(ctx, exec)
		if err != nil {
			return fmt.Errorf("failed to update todo: %w", err)
		}
	} else {
		// 存在しなければ挿入
		_, err = models.Todos.Insert(setter).Exec(ctx, exec)
		if err != nil {
			return fmt.Errorf("failed to insert todo: %w", err)
		}
	}

	return nil
}

// FindByID は、指定されたIDを持つTodoエンティティを検索します。
func (r *Repository) FindByID(ctx context.Context, id valueobject.TodoID) (*entity.Todo, error) {
	uid, err := uuid.Parse(id.String())
	if err != nil {
		return nil, fmt.Errorf("invalid todo id format: %w", err)
	}

	// bobが生成したFindTodoヘルパー関数を利用
	todoModel, err := models.FindTodo(ctx, r.getExecutor(ctx), uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewErr(errs.InternalCodeNotFound, fmt.Errorf("todo not found: %s", id.String()))
		}
		return nil, fmt.Errorf("failed to find todo by id: %w", err)
	}

	return toDomainTodo(todoModel)
}

// Delete は、指定されたIDを持つTodoエンティティを削除します。
func (r *Repository) Delete(ctx context.Context, id valueobject.TodoID) error {
	uid, err := uuid.Parse(id.String())
	if err != nil {
		return fmt.Errorf("invalid todo id format: %w", err)
	}

	_, err = models.Todos.Delete(
		dm.Where(models.Todos.Columns.ID.EQ(psql.Arg(uid))),
	).Exec(ctx, r.getExecutor(ctx))
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	return nil
}

// List は、すべてのTodoエンティティをsort_order順で取得します。
func (r *Repository) List(ctx context.Context) ([]*entity.Todo, error) {
	todoModels, err := models.Todos.Query(
		sm.OrderBy(models.Todos.Columns.SortOrder),
	).All(ctx, r.getExecutor(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}

	domainTodos := make([]*entity.Todo, len(todoModels))
	for i, m := range todoModels {
		domainTodos[i], err = toDomainTodo(m)
		if err != nil {
			return nil, fmt.Errorf("failed to convert todo model to domain model: %w", err)
		}
	}

	return domainTodos, nil
}

// toDomainTodo は Bob の models.Todo をドメインの todomodel.Todo に変換します。
func toDomainTodo(m *models.Todo) (*entity.Todo, error) {
	if m == nil {
		return nil, nil
	}

	description := m.Description.GetOrZero()

	var dueDate *time.Time
	if dd, ok := m.DueDate.Get(); ok {
		dueDate = &dd
	}

	return entity.Reconstruct(
		m.ID.String(),
		m.UserID.String(),
		m.Title,
		description,
		m.Status,
		dueDate,
		m.SortOrder,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

// toDBTodoSetter はドメインの todomodel.Todo を Bob の models.TodoSetter に変換します。
func toDBTodoSetter(td *entity.Todo) *models.TodoSetter {
	if td == nil {
		return nil
	}

	var description omitnull.Val[string]
	if td.Description() != "" {
		description = omitnull.From(td.Description())
	} else {
		description = omitnull.FromNull(null.Val[string]{})
	}

	// DueDate: nil のときは zero value (unset) のまま → INSERT/UPDATE に含めない
	// 値があるときだけセット、ClearDueDate は上位レイヤー(usecase)が nil を渡すことで表現
	var dueDate omitnull.Val[time.Time]
	if dd := td.DueDate(); dd != nil {
		dueDate = omitnull.From(*dd)
	}

	return &models.TodoSetter{
		ID:          omit.From(uuid.MustParse(td.ID().String())),
		UserID:      omit.From(uuid.MustParse(td.UserID().String())),
		Title:       omit.From(td.Title().String()),
		Description: description,
		Status:      omit.From(td.Status().String()),
		DueDate:     dueDate,
		SortOrder:   omit.From(td.SortOrder()),
		CreatedAt:   omit.From(td.CreatedAt()),
		UpdatedAt:   omit.From(td.UpdatedAt()),
	}
}

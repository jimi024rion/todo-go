package repository

import (
	"context"
	"database/sql" // sql.ErrNoRows のため
	"errors"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/aarondl/opt/null"

	"github.com/jimi024rion/todo-go/internal/domain/todo"
	"github.com/jimi024rion/todo-go/internal/domain/todo/entity"
	"github.com/jimi024rion/todo-go/models"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
)

type todoRepository struct {
	db bob.Executor
}

// NewTodoRepository creates a new todo repository.
func NewTodoRepository(db bob.Executor) todo.TodoRepository { // 型を変更
	return &todoRepository{db: db}
}

// toModel converts a domain entity.Todo to a Bob models.Todo
func toModel(e *entity.Todo) *models.Todo {
	m := &models.Todo{
		ID:        int32(e.ID),
		Title:     e.Title,
		Completed: e.Completed,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
	// entity.Description は string 型なので nil ではない
	// 空文字であれば Bob の null.Val を使う
	if e.Description == "" { // ここを修正
		m.Description = null.Val[string]{} // nil を表す null.Val
	} else {
		m.Description = null.From(e.Description) // 値を持つ null.Val
	}
	return m
}

// toEntity converts a Bob models.Todo to a domain entity.Todo
func toEntity(m *models.Todo) *entity.Todo {
	e := &entity.Todo{
		ID:        int64(m.ID),
		Title:     m.Title,
		Completed: m.Completed,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if !m.Description.IsNull() { // IsPresent の代わりに !IsNull() を使う
		e.Description = m.Description.MustGet() // *string ではなく string を取得
	} else {
		e.Description = "" // null の場合は空文字をセット
	}
	return e
}


func (r *todoRepository) Create(ctx context.Context, e *entity.Todo) (*entity.Todo, error) {
	// entity.Todo から models.TodoSetter を作成
	setter := &models.TodoSetter{
		Title:       omit.From(e.Title),
		Completed:   omit.From(e.Completed),
		CreatedAt:   omit.From(e.CreatedAt),
		UpdatedAt:   omit.From(e.UpdatedAt),
	}
	if e.Description != "" { // ここを修正 (nil ではなく空文字チェック)
		setter.Description = omitnull.From(e.Description) // string をそのままFrom
	} else {
		setter.Description = omitnull.Val[string]{} // 空文字ならDBにNULLを挿入
	}

	// 挿入を実行し、結果のモデルを取得
	m, err := models.Todos.Insert(
		setter,
		im.Returning("*"),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}

	// models.Todo から entity.Todo に変換して返す
	return toEntity(m), nil
}

func (r *todoRepository) FindByID(ctx context.Context, id int64) (*entity.Todo, error) {
	m, err := models.FindTodo(ctx, r.db, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return toEntity(m), nil
}

func (r *todoRepository) FindAll(ctx context.Context) ([]*entity.Todo, error) {
	ms, err := models.Todos.Query(
		sm.OrderBy(models.Todos.Columns.CreatedAt).Desc(),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	es := make([]*entity.Todo, len(ms))
	for i, m := range ms {
		es[i] = toEntity(m)
	}
	return es, nil
}

func (r *todoRepository) Update(ctx context.Context, e *entity.Todo) (*entity.Todo, error) {
	setter := models.TodoSetter{
		Title:       omit.From(e.Title),
		Completed:   omit.From(e.Completed),
		UpdatedAt:   omit.From(e.UpdatedAt),
	}
	if e.Description != "" { // ここを修正 (nil ではなく空文字チェック)
		setter.Description = omitnull.From(e.Description)
	} else {
		setter.Description = omitnull.Val[string]{} // 空文字ならDBにNULLを挿入
	}

	// IDでフィルタリングし、setterの内容で更新
	m, err := models.Todos.Update(
		setter.UpdateMod(),
		um.Where(models.Todos.Columns.ID.EQ(psql.Arg(e.ID))),
		um.Returning("*"),
	).One(ctx, r.db)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}

	return toEntity(m), nil
}

func (r *todoRepository) Delete(ctx context.Context, id int64) error {
	rowsAffected, err := models.Todos.Delete(
		dm.Where(models.Todos.Columns.ID.EQ(psql.Arg(int32(id)))),
	).Exec(ctx, r.db)

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Not found
	}

	return nil
}

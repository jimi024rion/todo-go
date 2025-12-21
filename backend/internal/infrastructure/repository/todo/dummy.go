package todo

import (
	"context"
	"fmt"
	"sync"

	todomodel "github.com/jimi024rion/todo-go/backend/internal/domain/model/todo"
	todorepository "github.com/jimi024rion/todo-go/backend/internal/domain/repository/todo"
)

var _ todorepository.TodoRepository = (*DummyRepository)(nil)

// DummyRepository は、テストやDI設定のためのTodoRepositoryのダミー実装です。
// データをインメモリのマップに保存します。
type DummyRepository struct {
	mu    sync.RWMutex
	todos map[todomodel.TodoID]*todomodel.Todo
}

// NewDummyRepository は新しいDummyRepositoryを生成します。
// この関数をwireでプロバイダーとして使用します。
func NewDummyRepository() *DummyRepository {
	return &DummyRepository{
		todos: make(map[todomodel.TodoID]*todomodel.Todo),
	}
}

// Save は、Todoエンティティをインメモリマップに保存します。
func (r *DummyRepository) Save(ctx context.Context, todo *todomodel.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.todos[todo.ID()] = todo
	fmt.Printf("DummyRepo: Saved todo %s\n", todo.ID())
	return nil
}

// FindByID は、インメモリマップから指定されたIDを持つTodoエンティティを検索します。
func (r *DummyRepository) FindByID(ctx context.Context, id todomodel.TodoID) (*todomodel.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if todo, ok := r.todos[id]; ok {
		fmt.Printf("DummyRepo: Found todo %s\n", id)
		return todo, nil
	}
	// 本来はインフラ層で定義したErrNotFoundを返すべき
	return nil, fmt.Errorf("todo with id %s not found in dummy repo", id)
}

// Delete は、インメモリマップから指定されたIDを持つTodoエンティティを削除します。
func (r *DummyRepository) Delete(ctx context.Context, id todomodel.TodoID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.todos[id]; !ok {
		return fmt.Errorf("todo with id %s not found in dummy repo", id)
	}
	delete(r.todos, id)
	fmt.Printf("DummyRepo: Deleted todo %s\n", id)
	return nil
}

// List は、インメモリマップに保存されているすべてのTodoエンティティのリストを取得します。
func (r *DummyRepository) List(ctx context.Context) ([]*todomodel.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*todomodel.Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		list = append(list, todo)
	}
	fmt.Printf("DummyRepo: Listed %d todos\n", len(list))
	return list, nil
}

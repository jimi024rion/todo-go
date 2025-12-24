// package entity は、Todoエンティティを定義します。
package entity

import (
	"time"

	"github.com/jimi024rion/todo-go/backend/internal/domain/todo/model/valueobject"
	uservo "github.com/jimi024rion/todo-go/backend/internal/domain/user/model/valueobject"
)

// --- エンティティ (Entity) ---

// Todoは、一つのタスクを表すエンティティです。
// フィールドはすべて非公開であり、公開されたメソッドを通じてのみ操作されます。
type Todo struct {
	id          valueobject.TodoID
	userID      uservo.UserID
	title       valueobject.Title
	description string
	status      valueobject.Status
	createdAt   time.Time
	updatedAt   time.Time
}

// --- ファクトリ (Factories) ---

// NewTodoは、新しいTodoエンティティを「新規作成」するためのファクトリ関数です。
// ID、ステータス、タイムスタンプは自動的に設定されます。
// 引数で受け取った値のビジネスルール検証もこの中で行います。
func NewTodo(userID uservo.UserID, title, description string) (*Todo, error) {
	// 値オブジェクトを生成することで、ビジネスルールを検証する
	validatedTitle, err := valueobject.NewTitle(title)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Todo{
		id:          valueobject.NewTodoID(),
		userID:      userID,
		title:       validatedTitle,
		description: description,
		status:      valueobject.StatusPending, // 新規作成時は必ず「未完了」
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// Reconstructは、データベースなどの永続化層から読み取ったデータを用いて、
// 既存のTodoエンティティをメモリ上に「再構築」するためのファクトリ関数です。
// この関数も、内部で値の検証を行います。
func Reconstruct(id, userID, title, description, status string, createdAt, updatedAt time.Time) (*Todo, error) {
	validatedID, err := valueobject.TodoIDFromString(id)
	if err != nil {
		return nil, err
	}
	validatedUserID, err := uservo.UserIDFromString(userID)
	if err != nil {
		return nil, err
	}
	validatedTitle, err := valueobject.NewTitle(title)
	if err != nil {
		return nil, err
	}
	validatedStatus, err := valueobject.NewStatus(status)
	if err != nil {
		return nil, err
	}

	return &Todo{
		id:          validatedID,
		userID:      validatedUserID,
		title:       validatedTitle,
		description: description,
		status:      validatedStatus,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}, nil
}

// --- メソッド (Methods) ---

// IDは、TodoのIDを返します。
func (t *Todo) ID() valueobject.TodoID {
	return t.id
}

// UserIDは、Todoの所有者であるユーザーのIDを返します。
func (t *Todo) UserID() uservo.UserID {
	return t.userID
}

// Titleは、Todoのタイトルを返します。
func (t *Todo) Title() valueobject.Title {
	return t.title
}

// Descriptionは、Todoの説明を返します。
func (t *Todo) Description() string {
	return t.description
}

// Statusは、Todoのステータスを返します。
func (t *Todo) Status() valueobject.Status {
	return t.status
}

// CreatedAtは、Todoが作成された日時を返します。
func (t *Todo) CreatedAt() time.Time {
	return t.createdAt
}

// UpdatedAtは、Todoが最後に更新された日時を返します。
func (t *Todo) UpdatedAt() time.Time {
	return t.updatedAt
}

// IsCompletedは、Todoが完了しているかどうかを返します。
func (t *Todo) IsCompleted() bool {
	return t.status == valueobject.StatusCompleted
}

// MarkAsCompletedは、Todoを「完了」状態に変更します。
// 冪等性を保つため、すでに完了状態の場合は何もせず、updatedAtも更新しません。
func (t *Todo) MarkAsCompleted() {
	if t.status == valueobject.StatusCompleted {
		return
	}
	t.status = valueobject.StatusCompleted
	t.updatedAt = time.Now()
}

// MarkAsInProgressは、Todoを「進行中」状態に変更します。
// 冪等性を保つため、すでに進行中状態の場合は何もせず、updatedAtも更新しません。
func (t *Todo) MarkAsInProgress() {
	if t.status == valueobject.StatusInProgress {
		return
	}
	t.status = valueobject.StatusInProgress
	t.updatedAt = time.Now()
}

// MarkAsPendingは、Todoを「未完了」状態に変更します。
// 冪等性を保つため、すでに未完了状態の場合は何もせず、updatedAtも更新しません。
func (t *Todo) MarkAsPending() {
	if t.status == valueobject.StatusPending {
		return
	}
	t.status = valueobject.StatusPending
	t.updatedAt = time.Now()
}

// ChangeTitleは、Todoのタイトルを変更します。
// 内部でビジネスルールを検証します。
func (t *Todo) ChangeTitle(newTitle string) error {
	validatedTitle, err := valueobject.NewTitle(newTitle)
	if err != nil {
		return err
	}
	if t.title == validatedTitle {
		return nil
	}
	t.title = validatedTitle
	t.updatedAt = time.Now()
	return nil
}

// ChangeDescriptionは、Todoの説明を変更します。
func (t *Todo) ChangeDescription(newDescription string) {
	if t.description == newDescription {
		return
	}
	t.description = newDescription
	t.updatedAt = time.Now()
}

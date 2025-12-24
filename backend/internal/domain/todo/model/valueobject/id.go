// package valueobject は、Todoに関連する値オブジェクトを定義します。
package valueobject

import (
	"fmt"

	"github.com/google/uuid"
)

// TodoIDは、Todoエンティティの一意な識別子を表す値オブジェクトです。
type TodoID uuid.UUID

// NewTodoIDは、新しいTodoIDを生成します。
func NewTodoID() TodoID {
	return TodoID(uuid.New())
}

// TodoIDFromStringは、文字列からTodoIDを生成します。
// 文字列が不正なUUID形式の場合はエラーを返します。
func TodoIDFromString(id string) (TodoID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return TodoID{}, fmt.Errorf("failed to parse todo id: %w", err)
	}
	return TodoID(uid), nil
}

// Stringは、TodoIDを文字列として返します。
func (id TodoID) String() string {
	return uuid.UUID(id).String()
}

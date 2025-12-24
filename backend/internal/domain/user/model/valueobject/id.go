package valueobject

import (
	"fmt"

	"github.com/google/uuid"
)

// UserID はユーザーの一意な識別子を表す値オブジェクトです。
type UserID uuid.UUID

// NewUserID は新しいUserIDを生成します。
func NewUserID() UserID {
	return UserID(uuid.New())
}

// UserIDFromString は文字列からUserIDを生成します。
// 文字列が不正なUUID形式の場合はエラーを返します。
func UserIDFromString(id string) (UserID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return UserID{}, fmt.Errorf("invalid user id format: %w", err)
	}
	return UserID(uid), nil
}

// String はUserIDを文字列として返します。
func (id UserID) String() string {
	return uuid.UUID(id).String()
}

// Value はプリミティブな値を返します。
func (id UserID) Value() string {
	return id.String()
}

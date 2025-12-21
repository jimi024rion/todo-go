// package todo は、Todoエンティティとそれに関連する値オブジェクトを定義します。
package todo

import (
	"errors"
	"fmt"
)

var (
	// ErrTitleIsEmpty は、タイトルが空の場合に返されるエラーです。
	ErrTitleIsEmpty = errors.New("title cannot be empty")
	// ErrTitleTooLong は、タイトルが長すぎる場合に返されるエラーです。
	ErrTitleTooLong = fmt.Errorf("title cannot be longer than %d characters", maxTitleLength)
)

const (
	// maxTitleLength はタイトルの最大長を定義します。
	maxTitleLength = 100
)

// Titleは、Todoのタイトルを表す値オブジェクトです。
// この型の値は、常にビジネスルール（空ではない、長すぎない）を満たしていることが保証されます。
type Title string

// NewTitleは、文字列からTitle値オブジェクトを生成します。
// ビジネスルール検証に失敗した場合はエラーを返します。
func NewTitle(title string) (Title, error) {
	if title == "" {
		return "", ErrTitleIsEmpty
	}
	if len(title) > maxTitleLength {
		return "", ErrTitleTooLong
	}
	return Title(title), nil
}

// Stringは、Titleを文字列として返します。
func (t Title) String() string {
	return string(t)
}

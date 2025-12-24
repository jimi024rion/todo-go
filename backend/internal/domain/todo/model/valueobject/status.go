// package valueobject は、Todoに関連する値オブジェクトを定義します。
package valueobject

import "errors"

// ErrInvalidStatus は、ステータスが無効な値の場合に返されるエラーです。
var ErrInvalidStatus = errors.New("invalid status")

// Statusは、Todoの進捗状態を表す値オブジェクトです。
type Status string

const (
	// StatusPendingは「未完了」の状態を表します。
	StatusPending Status = "pending"
	// StatusInProgressは「進行中」の状態を表します。
	StatusInProgress Status = "in_progress"
	// StatusCompletedは「完了」の状態を表します。
	StatusCompleted Status = "completed"
)

// NewStatusは、文字列からStatus値オブジェクトを生成します。
// 定義されていないステータス文字列の場合はエラーを返します。
func NewStatus(status string) (Status, error) {
	switch Status(status) {
	case StatusPending, StatusInProgress, StatusCompleted:
		return Status(status), nil
	default:
		return "", ErrInvalidStatus
	}
}

// Stringは、Statusを文字列として返します。
func (s Status) String() string {
	return string(s)
}

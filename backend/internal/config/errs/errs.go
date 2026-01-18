package errs

import (
	"github.com/cockroachdb/errors"
	pkgerrors "github.com/pkg/errors"
)

type ResultCode int

const (
	BadRequest            ResultCode = 4001
	IsUnauthorizedRequest ResultCode = 4002
	IsForbiddenRequest    ResultCode = 4003
	IsNotFoundRequest     ResultCode = 4004
	CodeUnknownErr        ResultCode = 9999
)

// Err はプロジェクト固有のエラー型です。
// 内部で cockroachdb/errors を利用してスタックトレースを保持します。
type Err struct {
	code ResultCode
	// cockroachdb/errors でラップされたエラーを保持
	err error
}

// NewErr は、エラーにスタックトレースとResultCodeを付与します。
func NewErr(code ResultCode, err error) error {
	return &Err{
		code: code,
		// errors.Wrapf を使ってスタックトレースをキャプチャする。
		// 第2引数のメッセージは空で良い。
		err: errors.Wrapf(err, ""),
	}
}

func (e *Err) Error() string {
	return e.err.Error()
}

// Unwrap は、errors.Unwrap, errors.Is, errors.As のために必要です。
func (e *Err) Unwrap() error {
	return e.err
}

func (e *Err) ResultCode() ResultCode {
	return e.code
}

// StackTrace は、pkgerrors.MarshalStack がスタックトレースを取得するために実装します。
// cockroachdb/errors はこのインターフェースを実装しているため、それを呼び出します。
func (e *Err) StackTrace() pkgerrors.StackTrace {
	type stackTracer interface {
		StackTrace() pkgerrors.StackTrace
	}
	var st stackTracer
	// 内部のエラーから StackTrace を取得しようと試みる
	if errors.As(e.err, &st) {
		return st.StackTrace()
	}
	return nil
}

// 以下、各種エラー判定用のヘルパー関数です。
func IsBadRequest(err error) bool {
	var e *Err
	if !errors.As(err, &e) {
		return false
	}
	return e.ResultCode() == BadRequest
}

func IsUnauthorized(err error) bool {
	var e *Err
	if !errors.As(err, &e) {
		return false
	}
	return e.ResultCode() == IsUnauthorizedRequest
}

func IsForbidden(err error) bool {
	var e *Err
	if !errors.As(err, &e) {
		return false
	}
	return e.ResultCode() == IsForbiddenRequest
}

func IsNotFound(err error) bool {
	var e *Err
	if !errors.As(err, &e) {
		return false
	}
	return e.ResultCode() == IsNotFoundRequest
}

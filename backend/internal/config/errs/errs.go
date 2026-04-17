package errs

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
)

// InternalCode はドメイン・ユースケース層で発生するエラーの細粒度コードです。
// usecase が NewErr に渡し、ResultCode へのマッピングに使われます。
type InternalCode int

const (
	InternalCodeUnknown InternalCode = 0
	// バリデーション系 (1xxx)
	InternalCodeTitleEmpty    InternalCode = 1001
	InternalCodeTitleTooLong  InternalCode = 1002
	InternalCodeInvalidStatus InternalCode = 1003
	InternalCodeInvalidID     InternalCode = 1004
	// リソース系 (2xxx)
	InternalCodeNotFound InternalCode = 2001
	// 認証認可系 (3xxx)
	InternalCodeUnauthorized InternalCode = 3001
	InternalCodeForbidden    InternalCode = 3002
	// システム系 (9xxx)
	InternalCodeInternal InternalCode = 9999
)

// ResultCode は API レスポンスの resultHeader.resultCode に設定される値です。
// HTTPStatus() と Message() でそれぞれ HTTP ステータスコードとメッセージを取得できます。
type ResultCode int

const (
	ResultOK              ResultCode = 0
	BadRequest            ResultCode = 4001
	IsUnauthorizedRequest ResultCode = 4002
	IsForbiddenRequest    ResultCode = 4003
	IsNotFoundRequest     ResultCode = 4004
	CodeUnknownErr        ResultCode = 9999
)

// internalToResult は InternalCode を ResultCode へマッピングします。
//
//nolint:gochecknoglobals
var internalToResult = map[InternalCode]ResultCode{
	InternalCodeTitleEmpty:    BadRequest,
	InternalCodeTitleTooLong:  BadRequest,
	InternalCodeInvalidStatus: BadRequest,
	InternalCodeInvalidID:     BadRequest,
	InternalCodeNotFound:      IsNotFoundRequest,
	InternalCodeUnauthorized:  IsUnauthorizedRequest,
	InternalCodeForbidden:     IsForbiddenRequest,
	InternalCodeInternal:      CodeUnknownErr,
}

// HTTPStatus は ResultCode に対応する HTTP ステータスコードを返します。
func (rc ResultCode) HTTPStatus() int {
	switch rc {
	case BadRequest:
		return http.StatusBadRequest
	case IsUnauthorizedRequest:
		return http.StatusUnauthorized
	case IsForbiddenRequest:
		return http.StatusForbidden
	case IsNotFoundRequest:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// Message は ResultCode に対応するメッセージを返します。
func (rc ResultCode) Message() string {
	switch rc {
	case BadRequest:
		return "Bad Request"
	case IsUnauthorizedRequest:
		return "Unauthorized"
	case IsForbiddenRequest:
		return "Forbidden"
	case IsNotFoundRequest:
		return "Not Found"
	default:
		return "Internal Server Error"
	}
}

// Err は InternalCode とスタックトレースを保持するカスタムエラー型です。
type Err struct {
	code    InternalCode
	err     error
	callers []uintptr
}

// NewErr は指定した InternalCode を持つ Err を生成します。
func NewErr(code InternalCode, err error) error {
	callers := make([]uintptr, 64)
	n := runtime.Callers(2, callers)
	return &Err{code, err, callers[:n]}
}

// formatFuncName フルパスの関数名をパッケージ名を含まない短い形式に整形する。
// e.g. "github.com/foo/bar/baz.MyFunc" -> "MyFunc"
func formatFuncName(fullName string) string {
	if i := strings.LastIndex(fullName, "/"); i != -1 {
		fullName = fullName[i+1:] // baz.MyFunc
	}
	if i := strings.Index(fullName, "."); i != -1 {
		fullName = fullName[i+1:] // MyFunc
	}
	return fullName
}

func MarshalStack(err error) any {
	var e *Err
	if !errors.As(err, &e) {
		return nil
	}
	frames := runtime.CallersFrames(e.callers)

	var stack []map[string]any
	for {
		frame, more := frames.Next()
		stack = append(stack, map[string]any{
			"func":   formatFuncName(frame.Function),
			"line":   frame.Line,
			"source": filepath.Base(frame.File),
		})

		if !more {
			break
		}
	}
	return stack
}

func (e *Err) Error() string {
	if e.err == nil {
		return "error is nil"
	}
	return fmt.Sprintf("%d: %s", e.code, e.err.Error())
}

func (e *Err) Unwrap() error {
	return e.err
}

// ResultCode は InternalCode を ResultCode にマッピングして返します。
func (e *Err) ResultCode() ResultCode {
	if e == nil {
		return ResultOK
	}
	if rc, ok := internalToResult[e.code]; ok {
		return rc
	}
	return CodeUnknownErr
}

func AsErr(err error) *Err {
	if err == nil {
		return nil
	}

	var e *Err
	if errors.As(err, &e) {
		return e
	}

	// 想定外エラーはラップ
	return NewErr(InternalCodeInternal, err).(*Err)
}

// ResultCodeFrom はエラーから ResultCode を取得します。
// *Err 型でない場合は CodeUnknownErr を返します。
func ResultCodeFrom(err error) ResultCode {
	var e *Err
	if errors.As(err, &e) {
		return e.ResultCode()
	}
	return CodeUnknownErr
}

func IsBadRequest(err error) bool {
	return ResultCodeFrom(err) == BadRequest
}

func IsUnauthorized(err error) bool {
	return ResultCodeFrom(err) == IsUnauthorizedRequest
}

func IsForbidden(err error) bool {
	return ResultCodeFrom(err) == IsForbiddenRequest
}

func IsNotFound(err error) bool {
	return ResultCodeFrom(err) == IsNotFoundRequest
}

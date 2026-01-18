package errs

import (
	"errors"
	"fmt"
	"runtime"
)

type ResultCode int

const (
	BadRequest            ResultCode = 4001
	IsUnauthorizedRequest ResultCode = 4002
	IsForbiddenRequest    ResultCode = 4003
	IsNotFoundRequest     ResultCode = 4004
	CodeUnknownErr        ResultCode = 9999
)

type Err struct {
	code    ResultCode
	err     error
	callers []uintptr
}

func NewErr(code ResultCode, err error) error {
	callers := make([]uintptr, 64)
	n := runtime.Callers(2, callers)
	return &Err{code, err, callers[:n]}
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

func (e *Err) StackFrames() []uintptr {
	return e.callers
}

func (e *Err) ResultCode() ResultCode {
	if e == nil {
		return 0
	}
	return e.code
}

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

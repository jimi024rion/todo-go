package response

import "github.com/jimi024rion/todo-go/backend/internal/config/errs"

// ResultHeader はすべてのレスポンスに含まれるヘッダーです。
type ResultHeader struct {
	ResultCode int `json:"resultCode" example:"0"`
} // @name ResultHeader

type Response[T any] struct {
	ResultHeader ResultHeader `json:"resultHeader"`
	ResultBody   *T           `json:"resultBody"`
}

func NewResponse[T any](resultCode int, resultBody *T) Response[T] {
	return Response[T]{
		ResultHeader: ResultHeader{ResultCode: resultCode},
		ResultBody:   resultBody,
	}
}

// OKHeader は正常系レスポンス用のヘッダーを生成します（resultCode=0）。
func OKHeader() ResultHeader {
	return ResultHeader{ResultCode: 0}
}

// FailNull は resultBody が null の異常系レスポンスを生成します。
// 返り値は swaggo アノテーションと同じ ErrorNullResponse 型です。
func FailNull(rc errs.ResultCode) ErrorNullResponse {
	return ErrorNullResponse{
		ResultHeader: ResultHeader{ResultCode: int(rc)},
		ResultBody:   nil,
	}
}

type ErrorNullResponse struct {
	ResultHeader ResultHeader `json:"resultHeader"`
	ResultBody   *struct{}    `json:"resultBody"`
} // @name ErrorNullResponse

package response

// ErrorResponse はAPIエラーレスポンスの共通フォーマット。
type ErrorResponse struct {
	Error string `json:"error" example:"invalid request body"`
}

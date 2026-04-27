package response

type ErrorResponse struct {
	Message string `json:"message" example:"Not Found"`
} // @name ErrorResponse

func Fail(message string) ErrorResponse {
	return ErrorResponse{Message: message}
}

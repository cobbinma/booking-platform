package handlers

type errorCodes string

const (
	InvalidRequest errorCodes = "INVALID_REQUEST"
	InternalError             = "INTERNAL_ERROR"
)

type errorResponse struct {
	Code    errorCodes `json:"code"`
	Message string     `json:"message"`
}

func newErrorResponse(code errorCodes, message string) *errorResponse {
	return &errorResponse{
		Code:    code,
		Message: message,
	}
}

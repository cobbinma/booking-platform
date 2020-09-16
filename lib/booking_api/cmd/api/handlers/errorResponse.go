package handlers

type errorCodes string

const (
	InvalidRequest   errorCodes = "INVALID_REQUEST"
	InternalError               = "INTERNAL_ERROR"
	NoAvailableSlots            = "NO_AVAILABLE_SLOTS"
	VenueNotGiven               = "VENUE_NOT_GIVEN"
	VenueNotFound               = "VENUE_NOT_FOUND"
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

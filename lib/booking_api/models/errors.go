package models

import "fmt"

var (
	ErrInvalidRequest   = fmt.Errorf("INVALID_REQUEST")
	ErrInternalError    = fmt.Errorf("INTERNAL_ERROR")
	ErrNoAvailableSlots = fmt.Errorf("NO_AVAILABLE_SLOTS")
)

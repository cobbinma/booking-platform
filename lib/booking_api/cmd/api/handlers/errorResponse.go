package handlers

import (
	"errors"
	"github.com/cobbinma/booking-platform/lib/booking_api/gateways/venueAPI"
	"github.com/cobbinma/booking-platform/lib/booking_api/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

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
	return &errorResponse{Code: code, Message: message}
}

func WriteError(c echo.Context, err error) error {
	if errors.Is(err, models.ErrInvalidRequest) {
		return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, "incorrect user request"))
	}
	if errors.Is(err, ErrVenueNotGiven) {
		return c.JSON(http.StatusBadRequest, newErrorResponse(VenueNotFound, "venue was not given"))
	}
	if venueAPI.ErrVenueNotFound(err) {
		return c.JSON(http.StatusBadRequest, newErrorResponse(VenueNotGiven, "venue could not be found"))
	}
	if errors.Is(err, models.ErrNoAvailableSlots) {
		return c.JSON(http.StatusInternalServerError, newErrorResponse(NoAvailableSlots, "no available slots"))
	}
	return c.JSON(http.StatusInternalServerError, newErrorResponse(InternalError, "could not create booking"))
}

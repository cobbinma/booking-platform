package handlers

import (
	"fmt"
	"github.com/cobbinma/booking/lib/booking_api/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (h *Handlers) GetBookingsByDate(c echo.Context) error {
	ctx := c.Request().Context()

	date, err := getDateFromRequest(c)
	if err != nil {
		logrus.Error(fmt.Errorf("%s : %w", "could not get date from request", err))
		message := "invalid date"
		return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, message))
	}

	bookings, err := h.repository.GetBookings(ctx, &models.BookingFilter{Date: &date})
	if err != nil {
		logrus.Error(fmt.Errorf("%s : %w", "could not get tables", err))
		message := "could not get tables"
		return c.JSON(http.StatusInternalServerError, newErrorResponse(InternalError, message))
	}

	return c.JSON(http.StatusOK, bookings)
}

func getDateFromRequest(c echo.Context) (models.Date, error) {
	date, err := models.DateFromString(c.Param("date"))
	if err != nil {
		return models.Date(time.Time{}), fmt.Errorf("%s : %w", "could not parse date", err)
	}

	return date, nil
}

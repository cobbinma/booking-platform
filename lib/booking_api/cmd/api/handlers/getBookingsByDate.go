package handlers

import (
	"fmt"
	"github.com/cobbinma/booking/lib/booking_api/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func GetBookingsByDate(repository models.Repository) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		date, err := getDateFromRequest(c)
		if err != nil {
			logrus.Error(fmt.Errorf("%s : %w", "could not get date from request", err))
			return WriteError(c, models.ErrInvalidRequest)
		}

		bookings, err := repository.GetBookings(ctx, models.BookingFilterWithDate(&date))
		if err != nil {
			logrus.Error(fmt.Errorf("%s : %w", "could not get bookings", err))
			return WriteError(c, models.ErrInternalError)
		}

		return c.JSON(http.StatusOK, bookings)
	}
}

func getDateFromRequest(c echo.Context) (models.Date, error) {
	date, err := models.DateFromString(c.Param("date"))
	if err != nil {
		return models.Date(time.Time{}), fmt.Errorf("%s : %w", "could not parse date", err)
	}

	return date, nil
}

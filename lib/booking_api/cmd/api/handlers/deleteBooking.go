package handlers

import (
	"fmt"
	"github.com/cobbinma/booking/lib/booking_api/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func DeleteBooking(repository models.Repository) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		id, err := getBookingIDFromRequest(c)
		if err != nil {
			logrus.Error(fmt.Errorf("%s : %w", "could not get id from request", err))
			message := "invalid id"
			return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, message))
		}

		err = repository.DeleteBookings(ctx, []int{id})
		if err != nil {
			logrus.Error(fmt.Errorf("%s : %w", "could not delete booking", err))
			message := "could not delete booking"
			return c.JSON(http.StatusInternalServerError, newErrorResponse(InternalError, message))
		}

		return c.NoContent(http.StatusOK)
	}
}

func getBookingIDFromRequest(c echo.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, fmt.Errorf("%s : %w", "could not parse id", err)
	}

	return id, nil
}

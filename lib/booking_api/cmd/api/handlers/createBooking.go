package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/cobbinma/booking/lib/booking_api/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func (h *Handlers) CreateBooking(c echo.Context) error {
	ctx := c.Request().Context()

	reqBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		logrus.Info(fmt.Errorf("%s : %w", "could not read request", err))
		return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, "incorrect user request"))
	}

	booking := models.NewBooking{}
	if err := json.Unmarshal(reqBody, &booking); err != nil {
		logrus.Info(fmt.Errorf("%s : %w", "could not unmarshall request", err))
		return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, "incorrect user request"))
	}

	if err := booking.Valid(ctx, h.tableClient); err != nil {
		logrus.Info(fmt.Errorf("%s : %w", "invalid request", err))
		message := fmt.Sprintf("incorrect user request : %s", err)
		return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, message))
	}

	bookings, err := h.repository.GetBookings(ctx, models.BookingFilterWithTableIDs([]models.TableID{booking.TableID}))
	if err != nil {
		logrus.Error(fmt.Errorf("%s : %w", "could not get bookings", err))
		message := "could not create booking"
		return c.JSON(http.StatusInternalServerError, newErrorResponse(InternalError, message))
	}

	for i := range bookings {
		if booking.StartsAt.Before(bookings[i].EndsAt) && bookings[i].StartsAt.Before(booking.EndsAt) {
			message := "incorrect user request : requested booking slot is not free"
			logrus.Info(fmt.Errorf(message))
			return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, message))
		}
	}

	if err := h.repository.CreateBooking(ctx, booking); err != nil {
		logrus.Error(fmt.Errorf("%s : %w", "could not create table", err))
		message := "could not create booking"
		return c.JSON(http.StatusInternalServerError, newErrorResponse(InternalError, message))
	}

	return c.NoContent(http.StatusCreated)
}

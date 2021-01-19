package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/cobbinma/booking-platform/lib/booking_api/models"
	"github.com/cobbinma/booking-platform/lib/booking_api/services"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func CreateBooking(repository models.Repository, tableClient models.TableClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		reqBody, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			logrus.Info(fmt.Errorf("%s : %w", "could not read request", err))
			return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, "incorrect user request"))
		}

		slot := models.Slot{}
		if err := json.Unmarshal(reqBody, &slot); err != nil {
			logrus.Info(fmt.Errorf("%s : %w", "could not unmarshall request", err))
			return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, "incorrect user request"))
		}

		service := services.NewCreateBookingService(repository, tableClient)
		booking, err := service.CreateBooking(ctx, slot)
		if err != nil {
			logrus.Info(fmt.Errorf("%s : %w", "service could not create booking", err))
			return WriteError(c, err)
		}

		return c.JSON(http.StatusCreated, booking)
	}
}

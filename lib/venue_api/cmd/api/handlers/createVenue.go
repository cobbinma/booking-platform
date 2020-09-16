package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/cobbinma/booking/lib/venue_api/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func CreateVenue(repository models.Repository) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		reqBody, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			logrus.Info(fmt.Errorf("%s : %w", "could not read request", err))
			return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, "incorrect user request"))
		}

		venue := models.VenueInput{}
		if err := json.Unmarshal(reqBody, &venue); err != nil {
			logrus.Info(fmt.Errorf("%s : %w", "could not unmarshall request", err))
			return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, "incorrect user request"))
		}

		if err := repository.CreateVenue(ctx, venue); err != nil {
			logrus.Info(fmt.Errorf("%s : %w", "could not create venue in repository", err))
			return c.JSON(http.StatusInternalServerError, newErrorResponse(InternalError, "internal error has occurred"))
		}

		return c.NoContent(http.StatusCreated)
	}
}

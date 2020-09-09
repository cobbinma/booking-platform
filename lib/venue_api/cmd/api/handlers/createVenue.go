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

func CreateVenue(c echo.Context) error {
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

	return nil
}

package handlers

import (
	"fmt"
	"github.com/cobbinma/booking-platform/lib/venue_api/models"
	"github.com/cobbinma/booking-platform/lib/venue_api/repositories/postgres"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func GetVenue(repository models.Repository) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		id, err := getVenueIDFromRequest(c)
		if err != nil {
			logrus.Error(fmt.Errorf("%s : %w", "could not get id from request", err))
			return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, "invalid id"))
		}

		venue, err := repository.GetVenue(ctx, id)
		if err != nil {
			m := "could not find venue"
			if postgres.ErrVenueNotFound(err) {
				logrus.Info(fmt.Errorf("%s : %w", m, err))
				return c.JSON(http.StatusNotFound, newErrorResponse(VenueNotFound, m))
			}
			logrus.Error(fmt.Errorf("%s : %w", m, err))
			return c.JSON(http.StatusInternalServerError, newErrorResponse(InternalError, m))
		}

		return c.JSON(http.StatusOK, venue)
	}
}

func getVenueIDFromRequest(c echo.Context) (models.VenueID, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, fmt.Errorf("%s : %w", "could not parse id", err)
	}

	venueID := models.NewVenueID(id)

	return venueID, nil
}

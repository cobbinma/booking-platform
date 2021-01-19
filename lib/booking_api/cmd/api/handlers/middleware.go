package handlers

import (
	"context"
	"fmt"
	"github.com/cobbinma/booking-platform/lib/booking_api/models"
	"github.com/labstack/echo/v4"
	"strconv"
)

var ErrVenueNotGiven = fmt.Errorf("VENUE_NOT_GIVEN")

func VenueMiddleware(next echo.HandlerFunc, client models.VenueClient) echo.HandlerFunc {
	return func(c echo.Context) error {
		venueID, err := strconv.Atoi(c.Param("venue_id"))
		if err != nil {
			return WriteError(c, ErrVenueNotGiven)
		}

		venue, err := client.GetVenue(c.Request().Context(), models.NewVenueID(venueID))
		if err != nil {
			return WriteError(c, err)
		}

		ctx := context.WithValue(c.Request().Context(), models.VenueCtxKey, *venue)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

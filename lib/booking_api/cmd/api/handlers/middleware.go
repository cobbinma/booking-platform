package handlers

import (
	"context"
	"github.com/cobbinma/booking/lib/booking_api/gateways/venueAPI"
	"github.com/cobbinma/booking/lib/booking_api/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func VenueMiddleware(next echo.HandlerFunc, client models.VenueClient) echo.HandlerFunc {
	return func(c echo.Context) error {
		venueID, err := strconv.Atoi(c.Param("venue_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, newErrorResponse(VenueNotGiven, "venue not given"))
		}

		venue, err := client.GetVenue(c.Request().Context(), models.NewVenueID(venueID))
		if err != nil {
			if venueAPI.ErrVenueNotFound(err) {
				return c.JSON(http.StatusBadRequest, newErrorResponse(VenueNotFound, "venue not found"))
			}
			return c.JSON(http.StatusInternalServerError, newErrorResponse(InternalError, "internal error"))
		}

		ctx := context.WithValue(c.Request().Context(), models.VenueCtxKey, *venue)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

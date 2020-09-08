package handlers

import (
	"context"
	"github.com/cobbinma/booking/lib/booking_api/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func VenueMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		venue := c.Param("venue_id")
		if venue == "" {
			return c.JSON(http.StatusBadRequest, newErrorResponse(VenueNotGiven, "venue not given"))
		}

		ctx := context.WithValue(c.Request().Context(), models.VenueCtxKey, models.NewVenueID(venue))
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

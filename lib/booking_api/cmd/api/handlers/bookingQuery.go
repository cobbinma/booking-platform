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

func BookingQuery(repository models.Repository, tableClient models.TableClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		reqBody, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			logrus.Info(fmt.Errorf("%s : %w", "could not read request", err))
			return WriteError(c, models.ErrInvalidRequest)
		}

		query := models.BookingQuery{}
		if err := json.Unmarshal(reqBody, &query); err != nil {
			logrus.Info(fmt.Errorf("%s : %w", "could not unmarshall request", err))
			return WriteError(c, models.ErrInvalidRequest)
		}

		if err := query.Valid(ctx); err != nil {
			logrus.Info(fmt.Errorf("%s : %w", "invalid request", err))
			return WriteError(c, models.ErrInvalidRequest)
		}

		tables, err := tableClient.GetTablesWithCapacity(ctx, query.People)
		if err != nil {
			logrus.Error(fmt.Errorf("%s : %w", "could not get tables with capacity", err))
			return WriteError(c, models.ErrInternalError)
		}

		if len(tables) == 0 {
			logrus.Info(fmt.Errorf("%s : %w", "no available tables", err))
			return WriteError(c, models.ErrNoAvailableSlots)
		}

		tableIDs := []models.TableID{}
		for i := range tables {
			tableIDs = append(tableIDs, tables[i].ID)
		}

		bookings, err := repository.GetBookings(
			ctx, models.BookingFilterWithTableIDs(tableIDs), models.BookingFilterWithDate(&query.Date))
		if err != nil {
			logrus.Error(fmt.Errorf("%s : %w", "could not get bookings", err))
			return WriteError(c, models.ErrInternalError)
		}

		for i := range tables {
			overlap := false
			for j := range bookings {
				if bookings[j].TableID == tables[i].ID && query.StartsAt.Before(bookings[j].EndsAt) && bookings[j].StartsAt.Before(query.EndsAt) {
					overlap = true
					break
				}
			}
			if overlap {
				continue
			}
			return c.JSON(http.StatusOK, models.NewBooking{
				CustomerID: query.CustomerID,
				TableID:    tables[i].ID,
				People:     query.People,
				Date:       query.Date,
				StartsAt:   query.StartsAt,
				EndsAt:     query.EndsAt,
			})
		}

		logrus.Info(fmt.Errorf("%s : %w", "no available slots", err))
		return WriteError(c, models.ErrNoAvailableSlots)
	}
}

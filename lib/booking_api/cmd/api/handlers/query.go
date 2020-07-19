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

func (h *Handlers) BookingQuery(c echo.Context) error {
	ctx := c.Request().Context()

	reqBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		logrus.Info(fmt.Errorf("%s : %w", "could not read request", err))
		return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, "incorrect user request"))
	}

	query := models.BookingQuery{}
	if err := json.Unmarshal(reqBody, &query); err != nil {
		logrus.Info(fmt.Errorf("%s : %w", "could not unmarshall request", err))
		return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, "incorrect user request"))
	}

	if err := query.Valid(); err != nil {
		logrus.Info(fmt.Errorf("%s : %w", "invalid request", err))
		message := fmt.Sprintf("incorrect user request : %s", err)
		return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, message))
	}

	tables, err := h.tableClient.GetTablesWithCapacity(ctx, query.People)
	if err != nil {
		logrus.Error(fmt.Errorf("%s : %w", "could not get tables with capacity", err))
		message := "could not get booking request"
		return c.JSON(http.StatusInternalServerError, newErrorResponse(InternalError, message))
	}

	if len(tables) == 0 {
		logrus.Info(fmt.Errorf("%s : %w", "no available tables", err))
		return c.JSON(http.StatusNotFound, newErrorResponse(NoAvailableSlots, "no available tables"))
	}

	tableIDs := []models.TableID{}
	for i := range tables {
		tableIDs = append(tableIDs, tables[i].ID)

	}

	bookings, err := h.repository.GetBookings(
		ctx, models.BookingFilterWithTableIDs(tableIDs), models.BookingFilterWithDate(&query.Date))
	if err != nil {
		logrus.Error(fmt.Errorf("%s : %w", "could not get bookings", err))
		message := "could not get booking request"
		return c.JSON(http.StatusInternalServerError, newErrorResponse(InternalError, message))
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
	return c.JSON(http.StatusNotFound, newErrorResponse(NoAvailableSlots, "no available slots"))
}

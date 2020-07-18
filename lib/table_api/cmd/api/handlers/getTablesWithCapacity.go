package handlers

import (
	"fmt"
	"github.com/cobbinma/booking/lib/table_api/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (h *Handlers) GetTablesWithCapacity(c echo.Context) error {
	ctx := c.Request().Context()

	capacity, err := getCapacityFromRequest(c)
	if err != nil {
		logrus.Error(fmt.Errorf("%s : %w", "could not get capacity from request", err))
		message := "invalid capacity"
		return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, message))
	}

	tables, err := h.repository.GetTables(ctx, models.NewTableFilter(capacity))
	if err != nil {
		logrus.Error(fmt.Errorf("%s : %w", "could not get tables", err))
		message := "could not get tables"
		return c.JSON(http.StatusInternalServerError, newErrorResponse(InternalError, message))
	}

	return c.JSON(http.StatusOK, tables)
}

func getCapacityFromRequest(c echo.Context) (models.Capacity, error) {
	amount, err := strconv.Atoi(c.Param("amount"))
	if err != nil {
		return 0, fmt.Errorf("%s : %w", "could not parse capacity amount", err)
	}

	capacity := models.NewCapacity(amount)
	err = capacity.Valid()
	if err != nil {
		return 0, fmt.Errorf("%s : %w", "could not validate capacity", err)
	}

	return capacity, nil
}

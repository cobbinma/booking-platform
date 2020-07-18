package handlers

import (
	"fmt"
	"github.com/cobbinma/booking/lib/table_api/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (h *Handlers) DeleteTable(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := getTableIDFromRequest(c)
	if err != nil {
		logrus.Error(fmt.Errorf("%s : %w", "could not get id from request", err))
		message := "invalid id"
		return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, message))
	}

	err = h.repository.DeleteTable(ctx, id)
	if err != nil {
		logrus.Error(fmt.Errorf("%s : %w", "could not delete table", err))
		message := "could not delete table"
		return c.JSON(http.StatusInternalServerError, newErrorResponse(InternalError, message))
	}

	return c.NoContent(http.StatusOK)
}

func getTableIDFromRequest(c echo.Context) (models.TableID, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, fmt.Errorf("%s : %w", "could not parse id", err)
	}

	tableID := models.NewTableID(id)
	err = tableID.Valid()
	if err != nil {
		return 0, fmt.Errorf("%s : %w", "could not validate id", err)
	}

	return tableID, nil
}

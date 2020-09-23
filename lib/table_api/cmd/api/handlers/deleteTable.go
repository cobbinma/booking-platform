package handlers

import (
	"fmt"
	"github.com/cobbinma/booking/lib/table_api/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

func DeleteTable(repository models.Repository) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		id, err := getTableIDFromRequest(c)
		if err != nil {
			logrus.Error(fmt.Errorf("%s : %w", "could not get id from request", err))
			message := "invalid id"
			return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, message))
		}

		err = repository.DeleteTable(ctx, id)
		if err != nil {
			logrus.Error(fmt.Errorf("%s : %w", "could not delete table", err))
			message := "could not delete table"
			return c.JSON(http.StatusInternalServerError, newErrorResponse(InternalError, message))
		}

		return c.NoContent(http.StatusOK)
	}
}

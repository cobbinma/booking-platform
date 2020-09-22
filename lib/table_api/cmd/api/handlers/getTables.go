package handlers

import (
	"fmt"
	"github.com/cobbinma/booking/lib/table_api/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

func GetTables(repository models.Repository) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		tables, err := repository.GetTables(ctx, nil)
		if err != nil {
			logrus.Error(fmt.Errorf("%s : %w", "could not get tables", err))
			message := "could not get tables"
			return c.JSON(http.StatusInternalServerError, newErrorResponse(InternalError, message))
		}

		return c.JSON(http.StatusOK, tables)
	}
}

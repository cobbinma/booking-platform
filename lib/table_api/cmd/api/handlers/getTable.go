package handlers

import (
	"fmt"
	"github.com/cobbinma/booking-platform/lib/table_api/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

func GetTable(repository models.Repository) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		id, err := getTableIDFromRequest(c)
		if err != nil {
			logrus.Error(fmt.Errorf("%s : %w", "could not get id from request", err))
			message := "invalid id"
			return c.JSON(http.StatusBadRequest, newErrorResponse(InvalidRequest, message))
		}

		tables, err := repository.GetTables(ctx, models.NewTableFilter(0, []models.TableID{id}))
		if err != nil {
			logrus.Error(fmt.Errorf("%s : %w", "could not get tables", err))
			message := "could not get table"
			return c.JSON(http.StatusInternalServerError, newErrorResponse(InternalError, message))
		}

		if tables == nil || len(tables) > 1 {
			logrus.Error(fmt.Errorf("%s : %w", "got incorrect repository response", err))
			message := "could not get table"
			return c.JSON(http.StatusInternalServerError, newErrorResponse(InternalError, message))
		}

		if len(tables) == 0 {
			logrus.Error(fmt.Errorf("%s : %w", "table id does not exist", err))
			message := "table does not exist"
			return c.JSON(http.StatusNotFound, newErrorResponse(InvalidRequest, message))
		}

		return c.JSON(http.StatusOK, tables[0])
	}
}

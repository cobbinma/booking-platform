package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handlers) GetTables(c echo.Context) error {
	ctx := c.Request().Context()

	tables, err := h.repository.GetTables(ctx, nil)
	if err != nil {
		logrus.Error(fmt.Errorf("%s : %w", "could not get tables", err))
		message := "could not get tables"
		return c.JSON(http.StatusInternalServerError, newErrorResponse(InternalError, message))
	}

	return c.JSON(http.StatusOK, tables)
}

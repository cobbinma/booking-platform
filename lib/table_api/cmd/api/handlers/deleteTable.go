package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
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

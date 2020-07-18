package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handlers) Health(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

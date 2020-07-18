package main

import (
	"github.com/cobbinma/booking/lib/table_api/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/healthz", health)

	e.Logger.Fatal(e.Start(config.Port()))
}

func health(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

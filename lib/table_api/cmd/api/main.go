package main

import (
	"context"
	"github.com/cobbinma/booking/lib/table_api/config"
	"github.com/cobbinma/booking/lib/table_api/repositories/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	dbClient, closeDB, err := postgres.NewDBClient()
	if err != nil {
		log.Fatalf("could not create database client : %v", err)
	}
	defer func() {
		if err := closeDB(); err != nil {
			log.Errorf("could not close database : %v", err)
		}
	}()

	repository := postgres.NewPostgres(dbClient)
	if err := repository.Migrate(context.Background()); err != nil {
		log.Fatalf("could not migrate : %v", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/healthz", health)

	e.Logger.Fatal(e.Start(config.Port()))
}

func health(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

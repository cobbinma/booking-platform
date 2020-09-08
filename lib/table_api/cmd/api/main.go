package main

import (
	"context"
	"github.com/cobbinma/booking/lib/table_api/cmd/api/handlers"
	"github.com/cobbinma/booking/lib/table_api/config"
	"github.com/cobbinma/booking/lib/table_api/repositories/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	dbClient, closeDB, err := postgres.NewDBClient()
	if err != nil {
		log.Fatal("could not create database client : ", err)
	}
	defer func() {
		if err := closeDB(); err != nil {
			log.Error("could not close database : ", err)
		}
	}()

	repository := postgres.NewPostgres(dbClient)
	if err := repository.Migrate(context.Background()); err != nil {
		log.Fatal("could not migrate : ", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	mw := handlers.VenueMiddleware

	h := handlers.NewHandlers(repository)

	e.GET("/healthz", h.Health)
	e.POST("/venues/:venue_id/tables", mw(h.CreateTable))
	e.GET("/venues/:venue_id/tables/:id", mw(h.GetTable))
	e.DELETE("/venues/:venue_id/tables/:id", mw(h.DeleteTable))
	e.GET("/venues/:venue_id/tables", mw(h.GetTables))
	e.GET("/venues/:venue_id/tables/capacity/:amount", mw(h.GetTablesWithCapacity))

	e.Logger.Fatal(e.Start(config.Port()))
}

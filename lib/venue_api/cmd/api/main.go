package main

import (
	"context"
	"github.com/cobbinma/booking/lib/venue_api/cmd/api/handlers"
	"github.com/cobbinma/booking/lib/venue_api/config"
	"github.com/cobbinma/booking/lib/venue_api/repositories/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	dbClient, closeDB, err := postgres.NewDBClient(config.PostgresURL())
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

	e.GET("/healthz", handlers.Health)
	e.POST("/venues", handlers.CreateVenue(repository))
	e.GET("/venues/:id", handlers.GetVenue(repository))
	e.DELETE("/venues/:id", handlers.DeleteVenue(repository))

	e.Logger.Fatal(e.Start(config.Port()))
}
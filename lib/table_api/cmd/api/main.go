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

	h := handlers.NewHandlers(repository)

	e.GET("/healthz", h.Health)
	e.PUT("/table", h.CreateTable)
	e.DELETE("/table/:id", h.DeleteTable)
	e.GET("/tables", h.GetTables)
	e.GET("/tables/capacity/:amount", h.GetTablesWithCapacity)

	e.Logger.Fatal(e.Start(config.Port()))
}

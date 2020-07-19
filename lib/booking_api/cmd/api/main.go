package main

import (
	"context"
	"github.com/cobbinma/booking/lib/booking_api/cmd/api/handlers"
	"github.com/cobbinma/booking/lib/booking_api/config"
	"github.com/cobbinma/booking/lib/booking_api/gateways/tableAPI"
	"github.com/cobbinma/booking/lib/booking_api/repositories/postgres"
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

	tableClient := tableAPI.NewTableAPI()

	e := echo.New()

	e.Use(middleware.Logger())

	h := handlers.NewHandlers(repository, tableClient)

	e.GET("/healthz", h.Health)
	e.PUT("/booking", h.CreateBooking)
	e.GET("/booking", h.BookingQuery)
	e.DELETE("/booking/:id", h.DeleteBooking)
	e.GET("/bookings/date/:date", h.GetBookingsByDate)

	e.Logger.Fatal(e.Start(config.Port()))
}
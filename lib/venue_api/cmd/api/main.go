package main

import (
	"context"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/api"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/models"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("could not listen : %s", err)
	}

	s := grpc.NewServer()
	api.RegisterVenueAPIServer(s, &Service{})

	log.Infof("starting gRPC listener on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %s", err)
	}
}

var _ api.VenueAPIServer = (*Service)(nil)

type Service struct{}

func (s Service) GetVenue(ctx context.Context, request *api.GetVenueRequest) (*models.Venue, error) {
	hours := []*models.OpeningHoursSpecification{
		{
			DayOfWeek: 1,
			Opens:     "10:00",
			Closes:    "19:00",
		},
		{
			DayOfWeek: 2,
			Opens:     "10:00",
			Closes:    "19:00",
		},
		{
			DayOfWeek: 3,
			Opens:     "10:00",
			Closes:    "19:00",
		},
		{
			DayOfWeek: 4,
			Opens:     "10:00",
			Closes:    "19:00",
		},
		{
			DayOfWeek: 5,
			Opens:     "10:00",
			Closes:    "19:00",
		},
		{
			DayOfWeek: 6,
			Opens:     "10:00",
			Closes:    "19:00",
		},
		{
			DayOfWeek: 7,
			Opens:     "10:00",
			Closes:    "19:00",
		},
	}
	return &models.Venue{
		Id:                  request.Id,
		Name:                "Hop and Vine",
		OpeningHours:        hours,
		SpecialOpeningHours: nil,
	}, nil
}

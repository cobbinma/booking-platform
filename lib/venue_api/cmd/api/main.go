package main

import (
	"context"
	"fmt"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/api"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/models"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("could not start logger : %s", err)
		os.Exit(-1)
	}
	defer logger.Sync()
	log := logger.Sugar()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("could not listen : %s", err)
	}

	s := grpc.NewServer(grpc_middleware.WithUnaryServerChain(
		grpc_zap.UnaryServerInterceptor(logger)))
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
			Closes:    "23:00",
		},
		{
			DayOfWeek: 7,
			Opens:     "10:00",
			Closes:    "20:00",
		},
	}
	return &models.Venue{
		Id:                  request.Id,
		Name:                "Hop and Vine",
		OpeningHours:        hours,
		SpecialOpeningHours: nil,
	}, nil
}

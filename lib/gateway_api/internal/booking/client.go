package booking

import (
	"context"
	"fmt"
	"github.com/cobbinma/booking-platform/lib/gateway_api/graph"
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/booking/api"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"time"
)

func NewBookingClient(url string, log *zap.SugaredLogger, token *oauth2.Token) (graph.BookingService, func(log *zap.SugaredLogger), error) {
	creds, err := credentials.NewClientTLSFromFile("localhost.crt", "localhost")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load credentials : %w", err)
	}

	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(oauth.NewOauthAccess(token)),
		grpc.WithTransportCredentials(creds),
	}
	conn, err := grpc.Dial(url, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect : %s", err)
	}

	return &bookingClient{
			client: api.NewBookingAPIClient(conn),
			log:    log,
		}, func(log *zap.SugaredLogger) {
			if err := conn.Close(); err != nil {
				log.Error("could not close connection : %s", err)
			}
		}, nil
}

type bookingClient struct {
	client api.BookingAPIClient
	log    *zap.SugaredLogger
}

func (b bookingClient) GetSlot(ctx context.Context, slot models.SlotInput) (*models.GetSlotResponse, error) {
	return &models.GetSlotResponse{
		Match: &models.Slot{
			VenueID:  slot.VenueID,
			Email:    slot.Email,
			People:   slot.People,
			StartsAt: slot.StartsAt,
			EndsAt:   slot.StartsAt.Add(time.Duration(slot.Duration) * time.Minute),
			Duration: slot.Duration,
		},
		OtherAvailableSlots: nil,
	}, nil
}

func (b bookingClient) CreateBooking(ctx context.Context, input models.BookingInput) (*models.Booking, error) {
	return &models.Booking{
		ID:       "5cbeadb9-b2b1-40ce-acbf-686f08f4e3af",
		VenueID:  input.VenueID,
		Email:    input.Email,
		People:   input.People,
		StartsAt: input.StartsAt,
		EndsAt:   input.StartsAt.Add(time.Duration(input.Duration) * time.Minute),
		Duration: input.Duration,
		TableID:  "6d3fe85d-a1cb-457c-bd53-48a40ee998e3",
	}, nil
}

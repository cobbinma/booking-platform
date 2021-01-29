package booking

import (
	"context"
	"fmt"
	"github.com/cobbinma/booking-platform/lib/gateway_api/graph"
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/booking/api"
	booking "github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/booking/models"
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
	resp, err := b.client.GetSlot(ctx, &booking.SlotInput{
		VenueId:  slot.VenueID,
		Email:    slot.Email,
		People:   (uint32)(slot.People),
		StartsAt: slot.StartsAt.Format(time.RFC3339),
		Duration: (uint32)(slot.Duration),
	})
	if err != nil {
		return nil, fmt.Errorf("could not get slot from booking api : %w", err)
	}

	var match *models.Slot
	if resp.Match != nil {
		startsAt, err := time.Parse(time.RFC3339, resp.Match.StartsAt)
		if err != nil {
			return nil, fmt.Errorf("could not parse start time : %w", err)
		}

		endsAt, err := time.Parse(time.RFC3339, resp.Match.EndsAt)
		if err != nil {
			return nil, fmt.Errorf("could not parse end time : %w", err)
		}

		match = &models.Slot{
			VenueID:  resp.Match.VenueId,
			Email:    resp.Match.Email,
			People:   (int)(resp.Match.People),
			StartsAt: startsAt,
			EndsAt:   endsAt,
			Duration: (int)(resp.Match.Duration),
		}
	}

	return &models.GetSlotResponse{
		Match:               match,
		OtherAvailableSlots: nil,
	}, nil
}

func (b bookingClient) CreateBooking(ctx context.Context, slot models.BookingInput) (*models.Booking, error) {
	resp, err := b.client.CreateBooking(ctx, &booking.SlotInput{
		VenueId:  slot.VenueID,
		Email:    slot.Email,
		People:   (uint32)(slot.People),
		StartsAt: slot.StartsAt.Format(time.RFC3339),
		Duration: (uint32)(slot.Duration),
	})
	if err != nil {
		return nil, fmt.Errorf("could not create booking in booking api : %w", err)
	}

	startsAt, err := time.Parse(time.RFC3339, resp.StartsAt)
	if err != nil {
		return nil, fmt.Errorf("could not parse start time : %w", err)
	}

	endsAt, err := time.Parse(time.RFC3339, resp.EndsAt)
	if err != nil {
		return nil, fmt.Errorf("could not parse end time : %w", err)
	}

	return &models.Booking{
		ID:       resp.Id,
		VenueID:  resp.VenueId,
		Email:    resp.Email,
		People:   (int)(resp.People),
		StartsAt: startsAt,
		EndsAt:   endsAt,
		Duration: (int)(resp.Duration),
		TableID:  resp.TableId,
	}, nil
}

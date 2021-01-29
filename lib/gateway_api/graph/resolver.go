package graph

import (
	"context"
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	venueService   VenueService
	userService    models.UserService
	bookingService BookingService
}

func NewResolver(userService models.UserService, venueService VenueService, bookingService BookingService) *Resolver {
	return &Resolver{userService: userService, venueService: venueService, bookingService: bookingService}
}

type VenueService interface {
	GetVenue(ctx context.Context, id string) (*models.Venue, error)
}

type BookingService interface {
	GetSlot(ctx context.Context, slot models.SlotInput) (*models.GetSlotResponse, error)
	CreateBooking(ctx context.Context, input models.BookingInput) (*models.Booking, error)
}

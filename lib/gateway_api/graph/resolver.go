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
	adminService   AdminService
}

func NewResolver(userService models.UserService, venueService VenueService, bookingService BookingService, adminService AdminService) *Resolver {
	return &Resolver{userService: userService, venueService: venueService, bookingService: bookingService, adminService: adminService}
}

//go:generate mockgen -package=mock_resolver -destination=./mock/graph.go -source=resolver.go
type VenueService interface {
	GetVenue(ctx context.Context, id string) (*models.Venue, error)
}

type BookingService interface {
	GetSlot(ctx context.Context, slot models.SlotInput) (*models.GetSlotResponse, error)
	CreateBooking(ctx context.Context, input models.BookingInput) (*models.Booking, error)
}

type AdminService interface {
	IsAdmin(ctx context.Context, input models.IsAdminInput) (bool, error)
}

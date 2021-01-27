package graph

import (
	"context"
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	venueService VenueService
	userService  models.UserService
}

func NewResolver(userService models.UserService, venueService VenueService) *Resolver {
	return &Resolver{userService: userService, venueService: venueService}
}

type VenueService interface {
	GetVenue(ctx context.Context, id string) (*models.Venue, error)
}

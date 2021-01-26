package graph

import (
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
	"github.com/cobbinma/booking-platform/lib/protobuf/autogen/lang/go/venue/api"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	venueService api.VenueAPIClient
	userService  models.UserService
}

func NewResolver(userService models.UserService, venueService api.VenueAPIClient) *Resolver {
	return &Resolver{userService: userService, venueService: venueService}
}

package graph

import "github.com/cobbinma/booking-platform/lib/gateway_api/models"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userService models.UserService
}

func NewResolver(userService models.UserService) *Resolver {
	return &Resolver{userService: userService}
}

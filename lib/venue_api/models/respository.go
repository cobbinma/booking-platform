package models

import "context"

type Repository interface {
	Migrate(ctx context.Context) error
	CreateVenue(ctx context.Context, venue VenueInput) error
}

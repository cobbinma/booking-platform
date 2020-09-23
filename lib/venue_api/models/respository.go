package models

import "context"

type Repository interface {
	Migrate(ctx context.Context, sourceURL string) error
	CreateVenue(ctx context.Context, venue VenueInput) (*Venue, error)
	GetVenue(ctx context.Context, id VenueID) (*Venue, error)
	DeleteVenue(ctx context.Context, id VenueID) error
}

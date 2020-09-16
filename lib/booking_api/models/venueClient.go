package models

import "context"

type VenueClient interface {
	GetVenue(ctx context.Context, id VenueID) (*Venue, error)
}

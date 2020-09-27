package models

import "context"

//go:generate mockgen -package=mock_models -destination=./mock/venueClient.go -source=venueClient.go
type VenueClient interface {
	GetVenue(ctx context.Context, id VenueID) (*Venue, error)
}

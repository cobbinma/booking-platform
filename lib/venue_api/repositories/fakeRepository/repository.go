package fakeRepository

import (
	"context"
	"database/sql"
	"github.com/cobbinma/booking/lib/venue_api/models"
)

type fakeRepository struct {
	venues map[models.VenueID]*models.Venue
}

func NewFakeRepository() models.Repository {
	return &fakeRepository{venues: make(map[models.VenueID]*models.Venue)}
}

func (f *fakeRepository) Migrate(ctx context.Context, sourceURL string) error {
	panic("implement me")
}

func (f *fakeRepository) CreateVenue(ctx context.Context, venue models.VenueInput) (*models.Venue, error) {
	id := models.VenueID(len(f.venues))
	newVenue := &models.Venue{
		ID:           id,
		Name:         venue.Name,
		OpeningHours: venue.OpeningHours,
	}
	f.venues[id] = newVenue
	return newVenue, nil
}

func (f *fakeRepository) GetVenue(ctx context.Context, id models.VenueID) (*models.Venue, error) {
	venue, ok := f.venues[id]
	if ok {
		return venue, nil
	}

	return nil, sql.ErrNoRows
}

func (f *fakeRepository) DeleteVenue(ctx context.Context, id models.VenueID) error {
	delete(f.venues, id)
	return nil
}

package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/cobbinma/booking-platform/lib/gateway_api/graph/generated"
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
)

func (r *mutationResolver) CreateSlot(ctx context.Context, input models.SlotInput) (*models.Slot, error) {
	starts, err := input.StartsAt.Time()
	if err != nil {
		return nil, fmt.Errorf("could not parse start time : %w", err)
	}

	return &models.Slot{
		VenueID:    input.VenueID,
		CustomerID: input.CustomerID,
		People:     input.People,
		Date:       input.Date,
		StartsAt:   input.StartsAt,
		EndsAt:     models.NewTimeOfDay(starts.Add(time.Minute * (time.Duration)(input.Duration))),
		Duration:   input.Duration,
	}, nil
}

func (r *mutationResolver) CreateBooking(ctx context.Context, input models.BookingInput) (*models.Booking, error) {
	user, err := r.userService.GetUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get user from context : %w", err)
	}

	if user.Email != input.CustomerID {
		return nil, fmt.Errorf("context email does not match given")
	}

	starts, err := input.StartsAt.Time()
	if err != nil {
		return nil, fmt.Errorf("could not parse start time : %w", err)
	}

	return &models.Booking{
		ID:         "5cbeadb9-b2b1-40ce-acbf-686f08f4e3af",
		VenueID:    input.VenueID,
		CustomerID: input.CustomerID,
		People:     input.People,
		Date:       input.Date,
		StartsAt:   input.StartsAt,
		EndsAt:     models.NewTimeOfDay(starts.Add(time.Minute * (time.Duration)(input.Duration))),
		Duration:   input.Duration,
		TableID:    "6d3fe85d-a1cb-457c-bd53-48a40ee998e3",
	}, nil
}

func (r *queryResolver) GetVenue(ctx context.Context, id string) (*models.Venue, error) {
	return &models.Venue{
		ID:   id,
		Name: "Hop and Vine",
		OpeningHours: []*models.OpeningHoursSpecification{
			{DayOfWeek: models.Monday, Opens: "10:00", Closes: "20:00"},
			{DayOfWeek: models.Tuesday, Opens: "10:00", Closes: "21:00"},
			{DayOfWeek: models.Wednesday, Opens: "10:00", Closes: "22:00"},
			{DayOfWeek: models.Thursday, Opens: "10:00", Closes: "22:00"},
			{DayOfWeek: models.Friday, Opens: "10:00", Closes: "22:00"},
			{DayOfWeek: models.Saturday, Opens: "10:00", Closes: "23:00"},
			{DayOfWeek: models.Sunday, Opens: "10:00", Closes: "23:00"},
		},
		SpecialOpeningHours: []*models.OpeningHoursSpecification{},
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

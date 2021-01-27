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

func (r *mutationResolver) CreateBooking(ctx context.Context, input models.BookingInput) (*models.Booking, error) {
	user, err := r.userService.GetUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get user from context : %w", err)
	}

	if user.Email != input.Email {
		return nil, fmt.Errorf("context email does not match given")
	}

	return &models.Booking{
		ID:       "5cbeadb9-b2b1-40ce-acbf-686f08f4e3af",
		VenueID:  input.VenueID,
		Email:    input.Email,
		People:   input.People,
		StartsAt: input.StartsAt,
		EndsAt:   input.StartsAt.Add(time.Duration(input.Duration) * time.Minute),
		Duration: input.Duration,
		TableID:  "6d3fe85d-a1cb-457c-bd53-48a40ee998e3",
	}, nil
}

func (r *queryResolver) GetVenue(ctx context.Context, id string) (*models.Venue, error) {
	return r.venueService.GetVenue(ctx, id)
}

func (r *queryResolver) GetSlot(ctx context.Context, input models.SlotInput) (*models.GetSlotResponse, error) {
	return &models.GetSlotResponse{
		Match: &models.Slot{
			VenueID:  input.VenueID,
			Email:    input.Email,
			People:   input.People,
			StartsAt: input.StartsAt,
			EndsAt:   input.StartsAt.Add(time.Duration(input.Duration) * time.Minute),
			Duration: input.Duration,
		},
		OtherAvailableSlots: nil,
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

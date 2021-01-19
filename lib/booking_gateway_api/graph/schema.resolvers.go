package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/cobbinma/booking-platform/lib/booking_gateway_api/graph/generated"
	"github.com/cobbinma/booking-platform/lib/booking_gateway_api/models"
)

func (r *mutationResolver) CreateSlot(ctx context.Context, input models.SlotInput) (*models.Slot, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateBooking(ctx context.Context, input models.BookingInput) (*models.Booking, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetVenue(ctx context.Context, id string) (*models.Venue, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

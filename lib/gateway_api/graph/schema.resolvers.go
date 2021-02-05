package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/cobbinma/booking-platform/lib/gateway_api/graph/generated"
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *mutationResolver) CreateBooking(ctx context.Context, input models.BookingInput) (*models.Booking, error) {
	user, err := r.userService.GetUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get user from context : %w", err)
	}

	if user.Email != input.Email {
		return nil, fmt.Errorf("context email does not match given")
	}

	return r.bookingService.CreateBooking(ctx, input)
}

func (r *queryResolver) GetVenue(ctx context.Context, id string) (*models.Venue, error) {
	return r.venueService.GetVenue(ctx, id)
}

func (r *queryResolver) GetSlot(ctx context.Context, input models.SlotInput) (*models.GetSlotResponse, error) {
	return r.bookingService.GetSlot(ctx, input)
}

func (r *queryResolver) IsAdmin(ctx context.Context, input models.IsAdminInput) (bool, error) {
	user, err := r.userService.GetUser(ctx)
	if err != nil {
		return false, status.Errorf(codes.Internal, "could not get user profile")
	}

	return r.customerService.IsAdmin(ctx, input.VenueID, user.Email)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

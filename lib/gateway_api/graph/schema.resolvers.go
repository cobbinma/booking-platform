package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/cobbinma/booking-platform/lib/gateway_api/graph/generated"
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *mutationResolver) CreateBooking(ctx context.Context, input models.BookingInput) (*models.Booking, error) {
	user, err := models.GetUserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get user profile : %w", err)
	}

	if user.Email != input.Email {
		if err := r.authIsAdmin(ctx, models.IsAdminInput{
			VenueID: &input.VenueID,
		}); err != nil {
			return nil, fmt.Errorf("context email does not match given : %w", err)
		}
	}

	return r.bookingService.CreateBooking(ctx, input)
}

func (r *mutationResolver) AddTable(ctx context.Context, input models.TableInput) (*models.Table, error) {
	if err := r.authIsAdmin(ctx, models.IsAdminInput{
		VenueID: &input.VenueID,
	}); err != nil {
		return nil, err
	}

	return r.venueService.AddTable(ctx, input)
}

func (r *mutationResolver) RemoveTable(ctx context.Context, input models.RemoveTableInput) (*models.Table, error) {
	if err := r.authIsAdmin(ctx, models.IsAdminInput{
		VenueID: &input.VenueID,
	}); err != nil {
		return nil, err
	}

	return r.venueService.RemoveTable(ctx, input)
}

func (r *mutationResolver) AddAdmin(ctx context.Context, input models.AdminInput) (string, error) {
	if err := r.authIsAdmin(ctx, models.IsAdminInput{
		VenueID: &input.VenueID,
	}); err != nil {
		return "", err
	}

	return r.venueService.AddAdmin(ctx, input)
}

func (r *mutationResolver) RemoveAdmin(ctx context.Context, input models.RemoveAdminInput) (string, error) {
	if err := r.authIsAdmin(ctx, models.IsAdminInput{
		VenueID: &input.VenueID,
	}); err != nil {
		return "", err
	}

	return r.venueService.RemoveAdmin(ctx, input)
}

func (r *mutationResolver) CancelBooking(ctx context.Context, input models.CancelBookingInput) (*models.Booking, error) {
	if input.VenueID == nil {
		return nil, fmt.Errorf("venue ID must be given")
	}
	if err := r.authIsAdmin(ctx, models.IsAdminInput{
		VenueID: input.VenueID,
	}); err != nil {
		return nil, err
	}

	return r.bookingService.CancelBooking(ctx, input)
}

func (r *mutationResolver) UpdateOpeningHours(ctx context.Context, input models.UpdateOpeningHoursInput) ([]*models.OpeningHoursSpecification, error) {
	r.log.Infof("updating opening hours : %v", input)
	if err := r.authIsAdmin(ctx, models.IsAdminInput{
		VenueID: &input.VenueID,
	}); err != nil {
		r.log.Errorf("user is not admin")
		return nil, err
	}

	hours, err := r.venueService.UpdateOpeningHours(ctx, input)
	if err != nil {
		r.log.Errorf("could not update opening hours : %s", err)
		return nil, fmt.Errorf("could not update opening hours : %w", err)
	}

	return hours, nil
}

func (r *mutationResolver) UpdateSpecialOpeningHours(ctx context.Context, input models.UpdateSpecialOpeningHoursInput) ([]*models.OpeningHoursSpecification, error) {
	if err := r.authIsAdmin(ctx, models.IsAdminInput{
		VenueID: &input.VenueID,
	}); err != nil {
		r.log.Errorf("user is not admin")
		return nil, err
	}

	hours, err := r.venueService.UpdateSpecialOpeningHours(ctx, input)
	if err != nil {
		r.log.Errorf("could not update special opening hours : %s", err)
		return nil, fmt.Errorf("could not update special opening hours : %w", err)
	}

	return hours, nil
}

func (r *queryResolver) GetVenue(ctx context.Context, filter models.VenueFilter) (*models.Venue, error) {
	if filter.ID == nil && filter.Slug == nil {
		return nil, fmt.Errorf("at least one field must not be nil on filter")
	}
	return r.venueService.GetVenue(ctx, filter)
}

func (r *queryResolver) GetSlot(ctx context.Context, input models.SlotInput) (*models.GetSlotResponse, error) {
	return r.bookingService.GetSlot(ctx, input)
}

func (r *queryResolver) IsAdmin(ctx context.Context, input models.IsAdminInput) (bool, error) {
	if input.VenueID == nil && input.Slug == nil {
		return false, fmt.Errorf("either venue id or slug must be given")
	}

	user, err := models.GetUserFromContext(ctx)
	if err != nil {
		return false, status.Errorf(codes.Internal, "could not get user profile : %s", err)
	}

	return r.venueService.IsAdmin(ctx, input, user.Email)
}

func (r *venueResolver) OpeningHoursSpecification(ctx context.Context, obj *models.Venue, date *time.Time) (*models.OpeningHoursSpecification, error) {
	if obj == nil || date == nil {
		return nil, nil
	}

	return r.venueService.OpeningHoursSpecification(ctx, obj.ID, *date)
}

func (r *venueResolver) Tables(ctx context.Context, obj *models.Venue) ([]*models.Table, error) {
	if err := r.authIsAdmin(ctx, models.IsAdminInput{
		VenueID: &obj.ID,
	}); err != nil {
		return nil, err
	}

	return r.venueService.GetTables(ctx, obj.ID)
}

func (r *venueResolver) Admins(ctx context.Context, obj *models.Venue) ([]string, error) {
	if err := r.authIsAdmin(ctx, models.IsAdminInput{
		VenueID: &obj.ID,
	}); err != nil {
		return nil, err
	}

	return r.venueService.GetAdmins(ctx, obj.ID)
}

func (r *venueResolver) Bookings(ctx context.Context, obj *models.Venue, filter *models.BookingsFilter, pageInfo *models.PageInfo) (*models.BookingsPage, error) {
	if err := r.authIsAdmin(ctx, models.IsAdminInput{
		VenueID: &obj.ID,
	}); err != nil {
		return nil, err
	}

	if filter == nil || pageInfo == nil {
		return nil, nil
	}

	if filter.VenueID != nil && *filter.VenueID != obj.ID {
		return nil, fmt.Errorf("cannot query bookings for a different venue")
	}

	return r.bookingService.Bookings(ctx, models.BookingsFilter{
		VenueID: &obj.ID,
		Date:    filter.Date,
	}, *pageInfo)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Venue returns generated.VenueResolver implementation.
func (r *Resolver) Venue() generated.VenueResolver { return &venueResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type venueResolver struct{ *Resolver }

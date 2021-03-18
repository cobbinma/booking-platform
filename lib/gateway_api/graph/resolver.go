package graph

import (
	"context"
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	log            *zap.SugaredLogger
	venueService   VenueService
	bookingService BookingService
	admins         *adminCache
}

func NewResolver(log *zap.SugaredLogger, venueService VenueService, bookingService BookingService) *Resolver {
	return &Resolver{
		log:            log,
		venueService:   venueService,
		bookingService: bookingService,
		admins:         newAdminCache(),
	}
}

//go:generate mockgen -package=mock_resolver -destination=./mock/graph.go -source=resolver.go
type VenueService interface {
	GetVenue(ctx context.Context, filter models.VenueFilter) (*models.Venue, error)
	OpeningHoursSpecification(ctx context.Context, venueID string, date time.Time) (*models.OpeningHoursSpecification, error)
	UpdateOpeningHours(ctx context.Context, input models.UpdateOpeningHoursInput) ([]*models.OpeningHoursSpecification, error)
	UpdateSpecialOpeningHours(ctx context.Context, input models.UpdateSpecialOpeningHoursInput) ([]*models.OpeningHoursSpecification, error)
	GetTables(ctx context.Context, venueID string) ([]*models.Table, error)
	AddTable(ctx context.Context, input models.TableInput) (*models.Table, error)
	RemoveTable(ctx context.Context, input models.RemoveTableInput) (*models.Table, error)
	IsAdmin(ctx context.Context, input models.IsAdminInput, email string) (bool, error)
	GetAdmins(ctx context.Context, venueID string) ([]string, error)
	AddAdmin(ctx context.Context, input models.AdminInput) (string, error)
	RemoveAdmin(ctx context.Context, input models.RemoveAdminInput) (string, error)
}

type BookingService interface {
	GetSlot(ctx context.Context, slot models.SlotInput) (*models.GetSlotResponse, error)
	CreateBooking(ctx context.Context, input models.BookingInput) (*models.Booking, error)
	Bookings(ctx context.Context, filter models.BookingsFilter, pageInfo models.PageInfo) (*models.BookingsPage, error)
	CancelBooking(ctx context.Context, input models.CancelBookingInput) (*models.Booking, error)
}

func (r *Resolver) authIsAdmin(ctx context.Context, input models.IsAdminInput) error {
	user, err := models.GetUserFromContext(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, "could not get user profile : %s", err)
	}

	if isAdmin := r.admins.getAdmin(user.Email); isAdmin {
		return nil
	}

	isAdmin, err := r.venueService.IsAdmin(ctx, input, user.Email)
	if err != nil {
		return status.Errorf(codes.Internal, "could not determine is user is admin : %s", err)
	}

	if isAdmin {
		r.admins.setAdmin(user.Email)
		return nil
	}

	return status.Errorf(codes.Unauthenticated, "user is not admin")
}

type adminCache struct {
	admins *cache.Cache
}

func newAdminCache() *adminCache {
	return &adminCache{admins: cache.New(5*time.Minute, 10*time.Minute)}
}

func (ac *adminCache) getAdmin(email string) bool {
	_, found := ac.admins.Get(email)
	return found
}

func (ac adminCache) setAdmin(email string) {
	ac.admins.Set(email, true, cache.DefaultExpiration)
}

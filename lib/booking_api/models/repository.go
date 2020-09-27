package models

import "context"

//go:generate mockgen -package=mock_models -destination=./mock/repository.go -source=repository.go
type Repository interface {
	Migrate(ctx context.Context, sourceURL string) error
	CreateBooking(ctx context.Context, booking Slot) (*Booking, error)
	GetBookings(ctx context.Context, options ...func(*BookingFilter) *BookingFilter) ([]Booking, error)
	DeleteBookings(ctx context.Context, id []int) error
}

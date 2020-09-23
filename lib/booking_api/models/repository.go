package models

import "context"

type Repository interface {
	Migrate(ctx context.Context, sourceURL string) error
	CreateBooking(ctx context.Context, booking NewBooking) (*Booking, error)
	GetBookings(ctx context.Context, options ...func(*BookingFilter) *BookingFilter) ([]Booking, error)
	DeleteBookings(ctx context.Context, id []int) error
}

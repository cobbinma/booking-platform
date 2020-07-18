package models

import "context"

type Repository interface {
	Migrate(ctx context.Context) error
	CreateBooking(ctx context.Context, booking NewBooking) error
	GetBookings(ctx context.Context, filter *BookingFilter) ([]Booking, error)
}

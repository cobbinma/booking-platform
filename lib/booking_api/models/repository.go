package models

import "context"

type Repository interface {
	Migrate(ctx context.Context) error
	CreateBooking(ctx context.Context, booking NewBooking) error
	GetBookings(ctx context.Context, options ...func(*BookingFilter) *BookingFilter) ([]Booking, error)
	DeleteBookings(ctx context.Context, id []int) error
}

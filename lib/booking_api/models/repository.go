package models

import "context"

type Repository interface {
	Migrate(ctx context.Context) error
	CreateBooking(ctx context.Context, booking NewBooking) error
}
